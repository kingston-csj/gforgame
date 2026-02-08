package handler

import (
	mysqldb "io/github/gforgame/db"
	playerdomain "io/github/gforgame/examples/domain/player"
	playerservice "io/github/gforgame/examples/service/player"
	"io/github/gforgame/examples/service/rank/container"
	"io/github/gforgame/examples/service/rank/model"
	"io/github/gforgame/protos"
)

type PlayerLevelRankHandler struct {
	BaseRankHandler
}

func NewPlayerLevelRankHandler() *PlayerLevelRankHandler {
	return &PlayerLevelRankHandler{BaseRankHandler: BaseRankHandler{rankContainer: container.NewConcurrentRankContainer(100)}}
}

func (p *PlayerLevelRankHandler) Init() {
	var players []playerdomain.Player
	err := mysqldb.Db.Where("level >?", 0).Order("level desc").Limit(50).Find(&players).Error
	if err != nil {
		panic("查询失败: " + err.Error())
	}
	for _, player := range players {
		p.rankContainer.Update(player.Id, NewPlayerLevelRank(player.Id, player.Level))
	}
}

func (p *PlayerLevelRankHandler) GetMyRankInfo(playerId string) *protos.RankInfo {
	player := playerservice.GetPlayerService().GetPlayer(playerId)
	rankInfo := &protos.RankInfo{
		Id: player.Id,
		Order: int32(p.QueryRankOrder(player.Id)),
		Value: int64(player.Level),
		SecondValue: 0,
		ExtraInfo:   "",
	}
	return rankInfo
}

func NewPlayerLevelRank(id string, level int32) *model.PlayerLevelRank {
	return &model.PlayerLevelRank{
		Id:    id,
		Level: level,
	}
}
