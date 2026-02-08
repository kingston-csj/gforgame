package rank

import (
	"io/github/gforgame/examples/service/player"
	"io/github/gforgame/examples/service/rank/handler"
	"io/github/gforgame/protos"
	"sync"
)

type RankType int

const (
	PlayerLevelRank    RankType = 1
	PlayerFightingRank RankType = 2
	PlayerArenaRank    RankType = 99
)

var (
	handlers map[RankType]handler.RankHandler = make(map[RankType]handler.RankHandler)
	once     sync.Once
	instance *RankService
)

// 排行榜模块
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

	playerArenaRank := handler.NewPlayerArenaRankHandler()
	playerArenaRank.Init()
	handlers[PlayerArenaRank] = playerArenaRank
}

func (rs *RankService) QueryRanks(rankType RankType, start int, end int) []protos.RankInfo {
	records := handlers[rankType].QueryRanks(start, end)
	rankInfos := make([]protos.RankInfo, 0)
	order := int32(start)
	for _, record := range records {
		vo := record.AsVo()
		vo.Name = player.GetPlayerService().GetPlayerProfileById(record.GetId()).Name
		vo.Order = order
		rankInfos = append(rankInfos, vo)
		order++
	}
	return rankInfos
}

func (rs *RankService) GetMyRankInfo(rankType RankType, playerId string) *protos.RankInfo {
	return handlers[rankType].GetMyRankInfo(playerId)
}
