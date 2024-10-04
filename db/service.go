package db

import (
	"io/github/gforgame/logger"
	"runtime"
	"strconv"
)

type AsyncDbService struct {
	run            int32
	workerCapacity int
	workers        []worker
}

func NewAsyncDbService() *AsyncDbService {
	capacity := max(4, runtime.NumCPU()/2)
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
	num, err := strconv.ParseInt((entity).GetId(), 10, 64)
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	index := num % int64(s.workerCapacity)
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
