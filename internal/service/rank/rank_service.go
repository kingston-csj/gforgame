package rank

import (
	"fmt"

	"github.com/forfun/gforgame/actor"
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/rank/container"
	"github.com/forfun/gforgame/internal/service/rank/handler"
)

type RankType int
const defaultRankCapacity = 100

const (
	PlayerLevelRank    RankType = 1
	PlayerFightingRank RankType = 2
	PlayerArenaRank    RankType = 99
)

var (
	handlers map[RankType]handler.RankHandler = make(map[RankType]handler.RankHandler)
	instance *RankService
)

// 排行榜模块
type RankService struct {
	playerRepo *playerrepo.PlayerRepository
	actorSystem *actor.ActorSystem
}

func NewRankService(playerRepo *playerrepo.PlayerRepository, actorSystem *actor.ActorSystem) *RankService {
	service := &RankService{
		playerRepo: playerRepo,
		actorSystem: actorSystem,
	}
	service.init()
	return service
}

func (rs *RankService) getRankActorRef(rankType RankType) *actor.ActorRef {
	path := actor.NewActorPath("game", "rank", fmt.Sprintf("%d", rankType))
	return rs.actorSystem.GetOrCreate (path, func() actor.Actor {
		return NewRankActor(rankType, defaultRankCapacity)
	})
}

func (rs *RankService) newRankContainer(rankType RankType) *container.ConcurrentRankContainer {
	ref := rs.getRankActorRef(rankType)
	return container.NewConcurrentRankContainer(ref)
}


func (rs *RankService) init() {
	levelContainer := rs.newRankContainer(PlayerLevelRank)
	playerLevelRank := handler.NewPlayerLevelRankHandler(rs.playerRepo, levelContainer)
	playerLevelRank.Init()
	handlers[PlayerLevelRank] = playerLevelRank

	playerFightingRank := handler.NewPlayerFightingRankHandler(rs.playerRepo, rs.newRankContainer(PlayerFightingRank))
	playerFightingRank.Init()
	handlers[PlayerFightingRank] = playerFightingRank

	playerArenaRank := handler.NewPlayerArenaRankHandler(rs.playerRepo, rs.newRankContainer(PlayerArenaRank))
	playerArenaRank.Init()
	handlers[PlayerArenaRank] = playerArenaRank
}

func (rs *RankService) QueryRanks(rankType RankType, start int, end int) []protos.RankInfo {
	records := handlers[rankType].QueryRanks(start, end)
	rankInfos := make([]protos.RankInfo, 0)
	order := int32(start)
	for _, record := range records {
		vo := record.AsVo()
		vo.Name = rs.playerRepo.GetPlayer(record.GetId()).Name
		vo.Order = order
		rankInfos = append(rankInfos, vo)
		order++
	}
	return rankInfos
}

func (rs *RankService) GetMyRankInfo(rankType RankType, playerId string) *protos.RankInfo {
	return handlers[rankType].GetMyRankInfo(playerId)
}
