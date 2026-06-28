package rank

import (
	"sync"

	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/internal/service/rank/handler"
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
	player *player.PlayerService
}

func NewRankService(playerService *player.PlayerService) *RankService {
	service := &RankService{
		player: playerService,
	}
	service.init()
	return service
}


func (rs *RankService) init() {
	playerLevelRank := handler.NewPlayerLevelRankHandler(rs.player)
	playerLevelRank.Init()
	handlers[PlayerLevelRank] = playerLevelRank

	playerFightingRank := handler.NewPlayerFightingRankHandler(rs.player)
	playerFightingRank.Init()
	handlers[PlayerFightingRank] = playerFightingRank

	playerArenaRank := handler.NewPlayerArenaRankHandler(rs.player)
	playerArenaRank.Init()
	handlers[PlayerArenaRank] = playerArenaRank
}

func (rs *RankService) QueryRanks(rankType RankType, start int, end int) []protos.RankInfo {
	records := handlers[rankType].QueryRanks(start, end)
	rankInfos := make([]protos.RankInfo, 0)
	order := int32(start)
	for _, record := range records {
		vo := record.AsVo()
		vo.Name = rs.player.GetPlayerProfileById(record.GetId()).Name
		vo.Order = order
		rankInfos = append(rankInfos, vo)
		order++
	}
	return rankInfos
}

func (rs *RankService) GetMyRankInfo(rankType RankType, playerId string) *protos.RankInfo {
	return handlers[rankType].GetMyRankInfo(playerId)
}
