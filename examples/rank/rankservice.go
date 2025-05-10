package rank

import (
	"io/github/gforgame/examples/rank/handler"
	"io/github/gforgame/examples/rank/model"
	"sync"
)

type RankType int

const (
	PlayerLevelRank    RankType = 1
	PlayerFightingRank RankType = 2
)

var (
	handlers map[RankType]handler.RankHandler = make(map[RankType]handler.RankHandler)
	once     sync.Once
	instance *RankService
)

type RankService struct {
}

func GetRankService() *RankService {
	once.Do(func() {
		instance = &RankService{}
		instance.init()
	})
	return instance
}

func (rs *RankService) init() {
	playerLevelRank := handler.NewPlayerLevelRankHandler()
	playerLevelRank.Init()
	handlers[PlayerLevelRank] = playerLevelRank

	playerFightingRank := handler.NewPlayerFightingRankHandler()
	playerFightingRank.Init()
	handlers[PlayerFightingRank] = playerFightingRank
}

func (rs *RankService) QueryRank(rankType RankType, start int, end int) []model.BaseRank {
	return handlers[rankType].QueryRanks(start, end)
}
