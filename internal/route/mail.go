package route

import (
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/mail"
)

type MailRoute struct {
	service    *mail.MailService
	playerrepo *playerrepo.PlayerRepository
}

func NewMailRoute(service *mail.MailService, playerRepo *playerrepo.PlayerRepository) *MailRoute {
	return &MailRoute{
		service:    service,
		playerrepo: playerRepo,
	}
}

func (c *MailRoute) ReqGetAllRewards(playerId string, index int32, msg *protos.ReqMailGetAllRewards) *protos.ResMailGetAllRewards {
	player := c.playerrepo.GetPlayer(playerId)
	rewardVos := c.service.TakeAllRewards(player)
	return &protos.ResMailGetAllRewards{
		Code:    0,
		Rewards: rewardVos,
	}
}

func (c *MailRoute) ReqDeleteAll(playerId string, index int32, msg *protos.ReqMailDeleteAll) *protos.ResMailDeleteAll {
	player := c.playerrepo.GetPlayer(playerId)
	removed := c.service.DeleteAll(player)
	return &protos.ResMailDeleteAll{
		Removed: removed,
	}
}

func (c *MailRoute) ReqGetReward(playerId string, index int32, msg *protos.ReqMailGetReward) *protos.ResMailGetReward {
	player := c.playerrepo.GetPlayer(playerId)
	code, rewardVos := c.service.TakeReward(player, msg.Id)
	if code != 0 {
		return &protos.ResMailGetReward{
			Code: int32(code),
		}
	}
	return &protos.ResMailGetReward{
		Rewards: rewardVos,
	}
}

func (c *MailRoute) ReqRead(playerId string, index int32, msg *protos.ReqMailRead) *protos.ResMailRead {
	player := c.playerrepo.GetPlayer(playerId)
	code := c.service.Read(player, msg.Id)
	return &protos.ResMailRead{
		Code: int32(code),
	}
}
