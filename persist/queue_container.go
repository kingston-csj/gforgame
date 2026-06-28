package persist

import (
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/forfun/gforgame/common/logger"
)

// ---------------------------------------------------
// QueueContainer 基于队列的持久化容器
// ---------------------------------------------------
type QueueContainer struct {
	name           string
	queue          chan string    // key 队列（同 key 只排一次）
	pending        sync.Map       // key -> latest snapshot Entity
	inQueue        sync.Map       // 并发安全 Set（是否已在队列中）
	savingStrategy SavingStrategy // 保存策略
	running        atomic.Bool    // 运行状态
	lastErrorTime  atomic.Int64   // 上次错误日志时间
	wg             sync.WaitGroup // 优雅关闭等待
}

func NewQueueContainer(name string, savingStrategy SavingStrategy) *QueueContainer {
	qc := &QueueContainer{
		name:           name,
		queue:          make(chan string, 1024*1024), // 大容量缓冲队列
		savingStrategy: savingStrategy,
	}
	qc.running.Store(true)

	// 启动后台协程
	qc.wg.Add(1)
	go qc.run()

	return qc
}

// Receive 接收实体
func (qc *QueueContainer) Receive(entity Entity) {
	if !qc.running.Load() {
		slog.Info("db closed, received entity", "key", entity.GetId())
		return
	}

	key := entity.GetId()
	snapshot, err := copyEntitySnapshot(entity)
	if err != nil {
		logger.ErrorNoStack("snapshot entity failed, key " + key + ", error: " + err.Error())
		return
	}
	qc.pending.Store(key, snapshot)

	// 去重：同 key 只排一次，队列内总是消费最新快照
	if _, loaded := qc.inQueue.LoadOrStore(key, struct{}{}); loaded {
		return
	}
	qc.queue <- key
}

// run 后台消费协程
func (qc *QueueContainer) run() {
	defer qc.wg.Done()

	for qc.running.Load() {
		select {
		case key, ok := <-qc.queue:
			if !ok {
				// channel关闭，退出goroutine
				return
			}
			qc.consumeKey(key)

		case <-time.After(1 * time.Second):
			// 每1秒轮询一次（同 queue.poll(1s)）
			continue
		}
	}
}

func (qc *QueueContainer) consumeKey(key string) {
	entityAny, ok := qc.pending.Load(key)
	if !ok {
		qc.inQueue.Delete(key)
		return
	}
	entity, ok := entityAny.(Entity)
	if !ok {
		qc.pending.Delete(key)
		qc.inQueue.Delete(key)
		return
	}

	err := qc.doSave(entity)
	if err != nil {
		// 失败重试：保持 inQueue 标记，重新入队同 key
		if qc.running.Load() {
			qc.queue <- key
		}
		now := time.Now().UnixMilli()
		last := qc.lastErrorTime.Load()
		if now-last > 5*60*1000 {
			qc.lastErrorTime.Store(now)
			logger.ErrorNoStack("save entity error, key " + key + ", error: " + err.Error())
		}
		return
	}

	latest, ok := qc.pending.Load(key)
	if ok && latest != entityAny {
		// 保存期间有新快照覆盖，继续消费最新版本
		if qc.running.Load() {
			qc.queue <- key
		}
		return
	}

	qc.pending.Delete(key)
	qc.inQueue.Delete(key)
}

// doSave 执行保存 + 异常处理（fatal 并发 map 写不会走到这里，必须靠快照规避）
func (qc *QueueContainer) doSave(entity Entity) (err error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic when save entity", "err", r, "key", entity.GetId())
			err = fmt.Errorf("panic when save entity: %v", r)
		}
	}()
	err = qc.savingStrategy.DoSave(entity)
	return
}

// ShutdownGraceful 优雅关闭
func (qc *QueueContainer) ShutdownGraceful() {
	qc.running.Store(false)
	close(qc.queue)
	qc.wg.Wait() // 等待消费协程退出

	// 关闭前把 pending 中剩余最新快照全部保存
	qc.pending.Range(func(k, v any) bool {
		key, _ := k.(string)
		entity, ok := v.(Entity)
		if !ok {
			return true
		}
		if err := qc.savingStrategy.DoSave(entity); err != nil {
			logger.ErrorNoStack("save entity on shutdown error, key " + key + ", error: " + err.Error())
		}
		return true
	})

	slog.Info("persist container shutdown gracefully", "name", qc.name)
}

// Size 当前队列大小
func (qc *QueueContainer) Size() int {
	size := 0
	qc.pending.Range(func(_, _ any) bool {
		size++
		return true
	})
	return size
}
