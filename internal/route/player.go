package route

import (
	"github.com/forfun/gforgame/common/eventbus"
	"github.com/forfun/gforgame/common/util/conv"
	"github.com/forfun/gforgame/internal/constants"
	"github.com/forfun/gforgame/internal/events"
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/service/player"

	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/network"
)

type PlayerRoute struct {
	service *player.PlayerService
	repo    *playerrepo.PlayerRepository
}

func NewPlayerRoute(service *player.PlayerService, repo *playerrepo.PlayerRepository) *PlayerRoute {
	return &PlayerRoute{
		service: service,
		repo:    repo,
	}
}

func (ps *PlayerRoute) ReqLogin(playerId string, s network.Session, index int32, msg *protos.ReqPlayerLogin) *protos.ResPlayerLogin {
	if conv.IsBlankString(msg.PlayerId) {
		s.Send(&protos.ResPlayerLogin{Code: constants.I18N_COMMON_ILLEGAL_PARAMS}, index)
		return &protos.ResPlayerLogin{Code: constants.I18N_COMMON_ILLEGAL_PARAMS}
	}
	return ps.service.DoLogin(msg.PlayerId, s, index)
}

func (ps *PlayerRoute) ReqLoadingFinish(playerId string, index int32, msg *protos.ReqPlayerLoadingFinish) {
	player := ps.repo.GetPlayer(playerId)
	eventbus.Default().Publish(events.PlayerLoadingFinish, player)
}

func (ps *PlayerRoute) ReqPlayerUpLevel(playerId string, index int32, msg *protos.ReqPlayerUpLevel) *protos.ResPlayerUpLevel {
	p := ps.repo.GetPlayer(playerId)
	return ps.service.DoUpLevel(p, msg.ToLevel)
}

func (ps *PlayerRoute) ReqPlayerUpStage(playerId string, index int32, msg *protos.ReqPlayerUpStage) *protos.ResPlayerUpStage {
	p := ps.repo.GetPlayer(playerId)
	return ps.service.DoUpStage(p)
}

func (ps *PlayerRoute) ReqPlayerRefreshScore(playerId string, index int32, msg *protos.ReqPlayerRefreshScore) *protos.ResPlayerRefreshScore {
	player := ps.repo.GetPlayer(playerId)
	player.ClientScore = msg.Score
	eventbus.Default().Publish(events.PlayerAttrChange, player)
	return &protos.ResPlayerRefreshScore{Code: 0}
}

func (ps *PlayerRoute) ReqEditClientData(playerId string, index int32, msg *protos.ReqEditClientData) *protos.ResEditClientData {
	player := ps.repo.GetPlayer(playerId)
	player.ClientData = msg.Data
	eventbus.Default().Publish(events.PlayerAttrChange, player)
	return &protos.ResEditClientData{Code: 0}
}
