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
	// 处理队列，失败的entity会被加入失败列表
	var failedEntities []struct {
		id     string
		entity Entity
	}
	for id, entity := range w.data {
			defer func() {
				w.mu.Unlock()
				if r := recover(); r != nil {
					var err error
					switch v := r.(type) {
					case error:
						err = v
					default:
						err = fmt.Errorf("%v", v)
					}
					logger.Error(err)
					// 单个entity处理失败，加入失败列表
					failedEntities = append(failedEntities, struct {
						id     string
						entity Entity
					}{id: id, entity: entity})
				}
			}()

		if entity.IsDeleted() {
			// 删除操作
			Db.Delete(entity)
		} else {
			// 保存操作
			entity.BeforeSave(nil)
			Db.Save(entity)
		}
		delete(w.data, id)
		atomic.AddInt32(&w.size, -1)
	}
	// 重新放入队列
	if len(failedEntities) > 0 {
		w.mu.Lock()
		defer w.mu.Unlock()
		for _, item := range failedEntities {
			w.data[item.id] = item.entity
		}
		// 恢复size（失败数）
		atomic.AddInt32(&w.size, int32(len(failedEntities)))
	}
}

func (w *worker) shutDown() {
	atomic.StoreInt32(&w.running, 0)
}

func (w *worker) QueueSize() int {
	return int(atomic.LoadInt32(&w.size))
}
