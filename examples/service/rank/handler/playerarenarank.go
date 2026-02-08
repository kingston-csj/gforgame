package handler

import (
	mysqldb "io/github/gforgame/db"
	playerdomain "io/github/gforgame/examples/domain/player"
	playerservice "io/github/gforgame/examples/service/player"
	"io/github/gforgame/examples/service/rank/container"
	"io/github/gforgame/examples/service/rank/model"
	"io/github/gforgame/protos"
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
