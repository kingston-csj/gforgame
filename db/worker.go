package db

import (
	"fmt"
	"io/github/gforgame/logger"
	"sync"
	"sync/atomic"
	"time"
)

type worker struct {
	data    map[string]Entity
	mu      sync.Mutex
	running int32
	size    int32
}

func (w *worker) run() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for atomic.LoadInt32(&w.running) == 1 {
		<-ticker.C
		w.processQueue()
	}
}

func (w *worker) addToQueue(entity Entity) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.data[entity.GetId()] = entity
	atomic.AddInt32(&w.size, 1)
}

func (w *worker) getEntity(id string) Entity {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.data[id]
}

func (w *worker) processQueue() {
	w.mu.Lock()
	if len(w.data) == 0 {
		w.mu.Unlock()
		return
	}

	// 1. 将数据转移到临时 map，并清空队列，以便释放锁
	pending := w.data
	w.data = make(map[string]Entity)
	atomic.StoreInt32(&w.size, 0)
	w.mu.Unlock()

	// 2. 处理数据（无锁状态）
	var failedEntities []Entity
	for _, entity := range pending {
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Error(fmt.Errorf("panic recovered: %v", r))
					failedEntities = append(failedEntities, entity)
				}
			}()

			if entity.IsDeleted() {
				Db.Delete(entity)
			} else {
				entity.BeforeSave(nil)
				Db.Save(entity)
			}
		}()
	}

	// 3. 将失败的任务重新放回队列
	if len(failedEntities) > 0 {
		w.mu.Lock()
		defer w.mu.Unlock()
		for _, entity := range failedEntities {
			id := entity.GetId()
			// 如果队列中已经有该 ID 的新数据，则放弃旧数据（避免版本回退）
			if _, exists := w.data[id]; !exists {
				w.data[id] = entity
				atomic.AddInt32(&w.size, 1)
			}
		}
	}
}

func (w *worker) shutDown() {
	atomic.StoreInt32(&w.running, 0)
}

func (w *worker) QueueSize() int {
	return int(atomic.LoadInt32(&w.size))
}
