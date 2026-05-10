package handler

import (
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	mysqldb "github.com/forfun/gforgame/internal/infra/persistence"
	"github.com/forfun/gforgame/internal/protos"
	playerservice "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/internal/service/rank/container"
	"github.com/forfun/gforgame/internal/service/rank/model"
)

type PlayerFightingRankHandler struct {
	BaseRankHandler
}

func NewPlayerFightingRankHandler() *PlayerFightingRankHandler {
	return &PlayerFightingRankHandler{BaseRankHandler: BaseRankHandler{rankContainer: container.NewConcurrentRankContainer(100)}}
}

func (p *PlayerFightingRankHandler) Init() {
	var players []playerdomain.Player
	err := mysqldb.Db.Where("fight >?", 0).Order("level desc").Limit(50).Find(&players).Error
	if err != nil {
		panic("查询失败: " + err.Error())
	}
	for _, player := range players {
		p.rankContainer.Update(player.Id, NewPlayerFightingRank(player.Id, player.Fight))
	}
}

func (p *PlayerFightingRankHandler) GetMyRankInfo(playerId string) *protos.RankInfo {
	player := playerservice.GetPlayerService().GetPlayer(playerId)
	rankInfo := &protos.RankInfo{
		Id: player.Id,
		Order: int32(p.QueryRankOrder(player.Id)),
		Value: int64(player.Fight),
		SecondValue: 0,
		ExtraInfo:   "",
	}
	return rankInfo
}

func NewPlayerFightingRank(id string, fighting int32) *model.PlayerFightingRank {
	return &model.PlayerFightingRank{
		Id:       id,
		Fighting: fighting,
	}
}
