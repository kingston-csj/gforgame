package route

import (
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/service/mixture"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"time"
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
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
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