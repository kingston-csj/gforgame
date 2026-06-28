package persist

import (
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/forfun/gforgame/common/logger"
)

// DelayContainer 延迟持久化容器
type DelayContainer struct {
	name         string
	delaySeconds int
	// key -> struct{}，表示该 key 是否已调度延迟任务（去重）
	pool sync.Map
	// key -> latest snapshot Entity
	pending        sync.Map
	savingStrategy SavingStrategy
	running        atomic.Bool    // 运行状态
	lastErrorTime  atomic.Int64   // 上次错误日志时间
	timer          *time.Timer    // 定时器（复用）
	wg             sync.WaitGroup // 优雅关闭
}

func NewDelayContainer(name string, delaySeconds int, savingStrategy SavingStrategy) *DelayContainer {
	dc := &DelayContainer{
		name:           name,
		delaySeconds:   delaySeconds,
		savingStrategy: savingStrategy,
	}
	dc.running.Store(true)
	return dc
}

// Receive 接收实体，延迟保存（去重 + 延迟执行）
func (dc *DelayContainer) Receive(entity Entity) {
	if !dc.running.Load() {
		logger.Info("db closed, received entity" + entity.GetId())
		return
	}

	key := entity.GetId()
	snapshot, err := copyEntitySnapshot(entity)
	if err != nil {
		logger.ErrorNoStack("snapshot entity failed, key " + key + ", error: " + err.Error())
		return
	}
	// 无论是否已调度，都要覆盖为最新快照
	dc.pending.Store(key, snapshot)

	// 去重：同 key 只调度一次延迟任务
	if _, loaded := dc.pool.LoadOrStore(key, struct{}{}); loaded {
		return
	}
	// 启动延迟任务（到期后落库最新快照）
	dc.wg.Add(1)
	go func() {
		defer dc.wg.Done()
		time.Sleep(time.Duration(dc.delaySeconds) * time.Second)
		dc.consumeLatest(key)
	}()
}

func (dc *DelayContainer) consumeLatest(key string) {
	snapAny, ok := dc.pending.Load(key)
	if !ok {
		dc.pool.Delete(key)
		return
	}
	entity, ok := snapAny.(Entity)
	if !ok {
		dc.pending.Delete(key)
		dc.pool.Delete(key)
		return
	}

	needRetry := false
	defer func() {
		if err := recover(); err != nil {
			logger.ErrorNoStack(err)
			needRetry = true
		}
		if needRetry {
			dc.pool.Delete(key)
			// 重试时继续走最新快照
			if latest, ok := dc.pending.Load(key); ok {
				if e, ok := latest.(Entity); ok {
					dc.Receive(e)
				}
			}
			return
		}

		latest, ok := dc.pending.Load(key)
		if ok && latest != snapAny {
			// 保存期间有更新，重新调度一次
			dc.pool.Delete(key)
			if e, ok := latest.(Entity); ok {
				dc.Receive(e)
			}
			return
		}

		dc.pending.Delete(key)
		dc.pool.Delete(key)
	}()

	err := dc.savingStrategy.DoSave(entity)
	if err != nil {
		needRetry = true
		now := time.Now().UnixMilli()
		last := dc.lastErrorTime.Load()
		if now-last > 5*60*1000 {
			dc.lastErrorTime.Store(now)
			logger.Error("save entity error, key ="+key, err)
		}
	}
}

// ShutdownGraceful 优雅关闭：立即执行所有未保存任务
func (dc *DelayContainer) ShutdownGraceful() {
	dc.running.Store(false)
	dc.wg.Wait() // 等待所有定时任务执行完

	// 把 pending 中剩下的全部立即保存
	dc.pending.Range(func(key, value any) bool {
		defer dc.pending.Delete(key)
		if entity, ok := value.(Entity); ok {
			if err := dc.savingStrategy.DoSave(entity); err != nil {
				logger.ErrorNoStack("save entity on shutdown error, key " + key.(string) + ", error: " + err.Error())
			}
		}
		return true
	})

	slog.Info("delay container shutdown gracefully", "name", dc.name)
}

// Size 获取当前待持久化数量
func (dc *DelayContainer) Size() int {
	size := 0
	dc.pool.Range(func(_, _ any) bool {
		size++
		return true
	})
	return size
}
