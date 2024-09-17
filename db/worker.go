package db

import (
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
	defer w.mu.Unlock()
	for id, entity := range w.data {
		if entity.IsDeleted() {
			// 删除操作
			Db.Delete(entity)
		} else {
			// 保存操作
			Db.Save(entity)
		}
		delete(w.data, id)
		atomic.AddInt32(&w.size, -1)
	}
}

func (w *worker) shutDown() {
	atomic.StoreInt32(&w.running, 0)
}

func (w *worker) queueSize() int {
	return int(atomic.LoadInt32(&w.size))
}
