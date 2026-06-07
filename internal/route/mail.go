package route

import (
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/mail"
	player "github.com/forfun/gforgame/internal/service/player"
)

type MailRoute struct {
	service *mail.MailService
	player  *player.PlayerService
}

func NewMailRoute(service *mail.MailService, playerService *player.PlayerService) *MailRoute {
	return &MailRoute{
		service: service,
		player:  playerService,
	}
}

func (c *MailRoute) ReqGetAllRewards(playerId string, index int32, msg *protos.ReqMailGetAllRewards) *protos.ResMailGetAllRewards {
	player := c.player.GetPlayer(playerId)
	rewardVos := c.service.TakeAllRewards(player)
	return &protos.ResMailGetAllRewards{
		Code: 0,
		Rewards: rewardVos,
	}
}

func (c *MailRoute) ReqDeleteAll(playerId string, index int32, msg *protos.ReqMailDeleteAll) *protos.ResMailDeleteAll {
	player := c.player.GetPlayer(playerId)
	removed := c.service.DeleteAll(player)
	return &protos.ResMailDeleteAll{
		Removed: removed,
	}
}

func (c *MailRoute) ReqGetReward(playerId string, index int32, msg *protos.ReqMailGetReward) *protos.ResMailGetReward {
	player := c.player.GetPlayer(playerId)
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
	player := c.player.GetPlayer(playerId)
	code := c.service.Read(player, msg.Id)
	return &protos.ResMailRead{
		Code: int32(code),
	}
}
