package persistence

import (
	"github.com/forfun/gforgame/examples/domain/player"
	"github.com/forfun/gforgame/persist"
)
type AsyncDBService struct {
	playerWorker persist.PersistContainer
	commonWorker persist.PersistContainer
}

func NewAsyncDbService() *AsyncDBService {
	return &AsyncDBService{
		playerWorker: persist.NewDelayContainer("player", 10, &EntitySavingStrategy{}),
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