package db

import (
	"hash/fnv"
	"runtime"
)

type AsyncDbService struct {
	run            int32
	workerCapacity int
	workers        []worker
}

func NewAsyncDbService() *AsyncDbService {
	capacity := max(4, runtime.NumCPU()/2)
	capacity = 1
	service := &AsyncDbService{
		workerCapacity: capacity,
		run:            1,
		workers:        make([]worker, capacity),
	}
	service.init()
	return service
}

func (s *AsyncDbService) init() {
	for i := range s.workers {
		s.workers[i] = worker{
			data:    make(map[string]Entity),
			running: 1,
		}
		go s.workers[i].run()
	}
}

func (s *AsyncDbService) SaveToDb(entity Entity) {
	if entity == nil {
		return
	}

	// 使用 fnv hash 计算字符串的哈希值
	h := fnv.New64a()
	h.Write([]byte(entity.GetId()))
	hash := h.Sum64()

	// 使用哈希值对worker容量取模来确定worker索引
	index := hash % uint64(s.workerCapacity)
	s.workers[index].addToQueue(entity)
}

func (s *AsyncDbService) DeleteFromDb(entity Entity) {
	entity.SetDeleted()
	s.SaveToDb(entity)
}

func (s *AsyncDbService) ShutDownGracefully() {
	for _, worker := range s.workers {
		worker.shutDown()
	}
}
