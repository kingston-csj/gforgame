package persist

import (
	"io/github/gforgame/common/logger"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

// DelayContainer 延迟持久化容器
type DelayContainer struct {
	name           string
	delaySeconds   int
	pool           sync.Map
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

	// 去重：已存在则不处理
	if _, exist := dc.pool.Load(key); exist {
		return
	}

	// 构建任务
	task := func() {
		needRetry := false
		defer func() {
			dc.pool.Delete(key)
			if err := recover(); err != nil {
				logger.ErrorNoStack(err)
				needRetry = true
			}
			if needRetry {
				dc.Receive(entity) // 失败或 panic 重试
			}
		}()

		// 执行保存
		err := dc.savingStrategy.DoSave(snapshot)
		if err != nil {
			needRetry = true

			// 5 分钟错误限流
			now := time.Now().UnixMilli()
			last := dc.lastErrorTime.Load()
			if now-last > 5*60*1000 {
				dc.lastErrorTime.Store(now)
				logger.Error("save entity error, key ="+entity.GetId(), err)
			}
			return
		}
	}

	// 存入池子（原子去重）
	_, loaded := dc.pool.LoadOrStore(key, task)
	if loaded {
		return
	}

	// 启动延迟任务
	dc.wg.Add(1)
	go func() {
		defer dc.wg.Done()
		time.Sleep(time.Duration(dc.delaySeconds) * time.Second)
		// if !dc.running.Load() {
		// 	return
		// }
		task()
	}()
}

// ShutdownGraceful 优雅关闭：立即执行所有未保存任务
func (dc *DelayContainer) ShutdownGraceful() {
	dc.running.Store(false)
	dc.wg.Wait() // 等待所有定时任务执行完

	// 把池子剩下的全部立即保存
	dc.pool.Range(func(key, task any) bool {
		defer dc.pool.Delete(key.(string))
		if t, ok := task.(func()); ok {
			t()
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
