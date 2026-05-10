package handler

import (
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	mysqldb "github.com/forfun/gforgame/internal/infra/persistence"
	"github.com/forfun/gforgame/internal/protos"
	playerservice "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/internal/service/rank/container"
	"github.com/forfun/gforgame/internal/service/rank/model"
)

type PlayerArenaRankHandler struct {
	BaseRankHandler
}

func NewPlayerArenaRankHandler() *PlayerArenaRankHandler {
	return &PlayerArenaRankHandler{BaseRankHandler: BaseRankHandler{rankContainer: container.NewConcurrentRankContainer(100)}}
}

func (p *PlayerArenaRankHandler) Init() {
	var players []playerdomain.Player
	err := mysqldb.Db.Where("level >?", 0).Order("level desc").Limit(50).Find(&players).Error
	if err != nil {
		panic("查询失败: " + err.Error())
	}
	for _, player := range players {
		p.rankContainer.Update(player.Id, NewPlayerLevelRank(player.Id, player.Level))
	}
}

func (p *PlayerArenaRankHandler) GetMyRankInfo(playerId string) *protos.RankInfo {
	player := playerservice.GetPlayerService().GetPlayer(playerId)
	rankInfo := &protos.RankInfo{
		Id: player.Id,
		Order: int32(p.QueryRankOrder(player.Id)),
		Value: int64(player.ArenaScore),
		SecondValue: 0,
		ExtraInfo:   "",
	}
	return rankInfo
}

func NewPlayerArenaRank(id string, level int32) *model.PlayerArenaRank {
	return &model.PlayerArenaRank{
		Id:    id,
		Score: level,
	}
}
