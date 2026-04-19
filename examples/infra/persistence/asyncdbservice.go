package persistence

import (
	"io/github/gforgame/examples/domain/player"
	"io/github/gforgame/persist"
)
type AsyncDBService struct {
	playerWorker *persist.QueueContainer
	commonWorker *persist.QueueContainer
}

func NewAsyncDbService() *AsyncDBService {
	return &AsyncDBService{
		playerWorker: persist.NewQueueContainer("player", &EntitySavingStrategy{}),
		commonWorker: persist.NewQueueContainer("common", &EntitySavingStrategy{}),
	}
}

func (s *AsyncDBService) SaveToDb(entity persist.Entity) {
	switch entity.(type) {
	case *player.Player:
		s.playerWorker.Receive(entity)
	default:
		s.commonWorker.Receive(entity)
	}
}

func (s *AsyncDBService) DeleteEntityFromDb(entity persist.Entity) {
	entity.SetDeleted()
	s.SaveToDb(entity)
}

func (s *AsyncDBService) Shutdown() {
	s.playerWorker.ShutdownGraceful()
	s.commonWorker.ShutdownGraceful()
}