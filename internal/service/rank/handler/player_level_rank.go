package handler

import (
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	mysqldb "github.com/forfun/gforgame/internal/infra/persistence"
	"github.com/forfun/gforgame/internal/protos"
	playerservice "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/internal/service/rank/container"
	"github.com/forfun/gforgame/internal/service/rank/model"
)

type PlayerLevelRankHandler struct {
	BaseRankHandler
	player *playerservice.PlayerService
}

func NewPlayerLevelRankHandler(player *playerservice.PlayerService) *PlayerLevelRankHandler {
	return &PlayerLevelRankHandler{
		BaseRankHandler: BaseRankHandler{rankContainer: container.NewConcurrentRankContainer(100)},
		player:          player,
	}
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
	player := p.player.GetPlayer(playerId)
	rankInfo := &protos.RankInfo{
		Id:          player.Id,
		Order:       int32(p.QueryRankOrder(player.Id)),
		Value:       int64(player.Level),
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
