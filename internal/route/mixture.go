package route

import (
	"time"

	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/mixture"
	"github.com/forfun/gforgame/network"

	player "github.com/forfun/gforgame/internal/service/player"
)

type MixtureRoute struct {
	service *mixture.MixtureService
	player  *player.PlayerService
}

func NewMixtureRoute(service *mixture.MixtureService, playerService *player.PlayerService) *MixtureRoute {
	return &MixtureRoute{
		service: service,
		player:  playerService,
	}
}

func (ps *MixtureRoute) ReqIdleViewReward(s *network.Session, index int32, msg *protos.ReqIdleViewReward) *protos.ResIdleViewReward {
	return &protos.ResIdleViewReward{
		Code: 0,
	}
}

func (ps *MixtureRoute) ReqClientUploadEvent(s *network.Session, index int32, msg *protos.ReqClientUploadEvent) *protos.ResClientUploadEvent {
	player := ps.player.GetPlayerBySession(s)
	ps.service.OnClientUploadEvent(player, msg.Type)
	return &protos.ResClientUploadEvent{
		Code: 0,
	}
}

func (c *MixtureRoute) ReqHeartBeat(s *network.Session, index int32, msg *protos.ReqHeartBeat) *protos.ResHeartBeat{
	return &protos.ResHeartBeat{
		Index: msg.Index,
		Code:  0,
	}
}

func (c *MixtureRoute) ReqGetServerTime(s *network.Session, index int32, msg *protos.ReqGetServerTime) *protos.ResGetServerTime{
	return &protos.ResGetServerTime{
		ServerTime: time.Now().Unix(),
		Code:  0,
	}
}