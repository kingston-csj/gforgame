package route

import (
	"time"

	"github.com/forfun/gforgame/examples/protos"
	"github.com/forfun/gforgame/examples/service/mixture"
	playerservice "github.com/forfun/gforgame/examples/service/player"
	"github.com/forfun/gforgame/network"
)

type MixtureRoute struct {
	network.Base
	service *mixture.MixtureService
}

func NewMixtureRoute() *MixtureRoute {
	return &MixtureRoute{
		service: mixture.GetMixtureService(),
	}
}

func (ps *MixtureRoute) Init() {
}

func (ps *MixtureRoute) ReqIdleViewReward(s *network.Session, index int32, msg *protos.ReqIdleViewReward) *protos.ResIdleViewReward {
	return &protos.ResIdleViewReward{
		Code: 0,
	}
}

func (ps *MixtureRoute) ReqClientUploadEvent(s *network.Session, index int32, msg *protos.ReqClientUploadEvent) *protos.ResClientUploadEvent {
	player := playerservice.GetPlayerService().GetPlayerBySession(s)
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