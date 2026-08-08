package route

import (
	"time"

	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/mixture"
)

type MixtureRoute struct {
	service    *mixture.MixtureService
	playerrepo *playerrepo.PlayerRepository
}

func NewMixtureRoute(service *mixture.MixtureService, playerRepo *playerrepo.PlayerRepository) *MixtureRoute {
	return &MixtureRoute{
		service:    service,
		playerrepo: playerRepo,
	}
}

func (ps *MixtureRoute) ReqIdleViewReward(playerId string, index int32, msg *protos.ReqIdleViewReward) *protos.ResIdleViewReward {
	return &protos.ResIdleViewReward{
		Code: 0,
	}
}

func (ps *MixtureRoute) ReqClientUploadEvent(playerId string, index int32, msg *protos.ReqClientUploadEvent) *protos.ResClientUploadEvent {
	player := ps.playerrepo.GetPlayer(playerId)
	ps.service.OnClientUploadEvent(player, msg.Type)
	return &protos.ResClientUploadEvent{
		Code: 0,
	}
}

func (c *MixtureRoute) ReqHeartBeat(playerId string, index int32, msg *protos.ReqHeartBeat) *protos.ResHeartBeat {
	return &protos.ResHeartBeat{
		Index: msg.Index,
		Code:  0,
	}
}

func (c *MixtureRoute) ReqGetServerTime(playerId string, index int32, msg *protos.ReqGetServerTime) *protos.ResGetServerTime {
	return &protos.ResGetServerTime{
		ServerTime: time.Now().Unix(),
		Code:       0,
	}
}
