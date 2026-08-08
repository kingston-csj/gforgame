package persistence

import (
	"runtime"

	"github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/persist"
)

type AsyncDBService struct {
	playerWorker persist.PersistContainer
	commonWorker persist.PersistContainer
}

func NewAsyncDbService() *AsyncDBService {
	// 系统核心数
	coreNum := runtime.NumCPU()
	playerWorkers := make([]persist.PersistContainer, 0, coreNum)
	entitySavingStrategy := &EntitySavingStrategy{}
	// 玩家数量多, 采用容器组
	// 数据变动频繁,采用delay方式,避免数据库压力过大
	for i := 0; i < coreNum; i++ {
		playerWorkers = append(playerWorkers, persist.NewDelayContainer("player", 3, entitySavingStrategy))
	}
	return &AsyncDBService{
		playerWorker: persist.NewPersistContainerGroup("player-group", playerWorkers),
		commonWorker: persist.NewQueueContainer("common", entitySavingStrategy),
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
