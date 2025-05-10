package handler

import (
	mysqldb "io/github/gforgame/db"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/rank/container"
	"io/github/gforgame/examples/rank/model"
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

func NewPlayerFightingRank(id string, fighting int32) *model.PlayerFightingRank {
	return &model.PlayerFightingRank{
		Id:       id,
		Fighting: fighting,
	}
}
