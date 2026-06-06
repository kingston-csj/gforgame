package route

import (
	"github.com/forfun/gforgame/common/util/conv"
	"github.com/forfun/gforgame/internal/constants"
	"github.com/forfun/gforgame/internal/context"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/service/player"

	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/network"
)

type PlayerRoute struct {
	service *player.PlayerService
}

func NewPlayerRoute(service *player.PlayerService) *PlayerRoute {
	return &PlayerRoute{
		service: service,
	}
}

func (ps *PlayerRoute) ReqLogin(s *network.Session, index int32, msg *protos.ReqPlayerLogin) {
	if conv.IsBlankString(msg.PlayerId) {
		s.Send(&protos.ResPlayerLogin{Code: constants.I18N_COMMON_ILLEGAL_PARAMS}, index)
		return
	}
	ps.service.DoLogin(msg.PlayerId, s, index)
}

func (ps *PlayerRoute) ReqLoadingFinish(s *network.Session, index int32, msg *protos.ReqPlayerLoadingFinish) {
	player := ps.service.GetPlayerBySession(s)
	context.EventBus.Publish(events.PlayerLoadingFinish, player)
}


func (ps *PlayerRoute) ReqPlayerUpLevel(s *network.Session, index int32, msg *protos.ReqPlayerUpLevel) *protos.ResPlayerUpLevel {
	p := ps.service.GetPlayerBySession(s)
	return ps.service.DoUpLevel(p, msg.ToLevel)
}

func (ps *PlayerRoute) ReqPlayerUpStage(s *network.Session, index int32, msg *protos.ReqPlayerUpStage) *protos.ResPlayerUpStage {
	p := ps.service.GetPlayerBySession(s)
	return ps.service.DoUpStage(p)
}

func (ps *PlayerRoute) ReqPlayerRefreshScore(s *network.Session, index int32, msg *protos.ReqPlayerRefreshScore) *protos.ResPlayerRefreshScore {
	player := ps.service.GetPlayerBySession(s)
	player.ClientScore = msg.Score
	context.EventBus.Publish(events.PlayerAttrChange, player)
	return &protos.ResPlayerRefreshScore{Code: 0}
}

func (ps *PlayerRoute) ReqEditClientData	(s *network.Session, index int32, msg *protos.ReqEditClientData) *protos.ResEditClientData {
	player := ps.service.GetPlayerBySession(s)
	player.ClientData = msg.Data
	context.EventBus.Publish(events.PlayerAttrChange, player)
	return &protos.ResEditClientData{Code: 0}
}