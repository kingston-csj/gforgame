package route

import (
	"github.com/forfun/gforgame/internal/context"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/mail"
	player "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/network"
)

type MailRoute struct {
	network.Base
	service *mail.MailService
	player  *player.PlayerService
}

func NewMailRoute(service *mail.MailService, playerService *player.PlayerService) *MailRoute {
	return &MailRoute{
		service: service,
		player:  playerService,
	}
}

func (c *MailRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		c.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (c *MailRoute) OnPlayerLogin(player *playerdomain.Player) {
	c.service.CheckMailsOnLogin(player)
}

func (c *MailRoute) ReqGetAllRewards(s *network.Session, index int32, msg *protos.ReqMailGetAllRewards) *protos.ResMailGetAllRewards {
	player := c.player.GetPlayerBySession(s)
	rewardVos := c.service.TakeAllRewards(player)
	return &protos.ResMailGetAllRewards{
		Code: 0,
		Rewards: rewardVos,
	}
}

func (c *MailRoute) ReqDeleteAll(s *network.Session, index int32, msg *protos.ReqMailDeleteAll) *protos.ResMailDeleteAll {
	player := c.player.GetPlayerBySession(s)
	removed := c.service.DeleteAll(player)
	return &protos.ResMailDeleteAll{
		Removed: removed,
	}
}

func (c *MailRoute) ReqGetReward(s *network.Session, index int32, msg *protos.ReqMailGetReward) *protos.ResMailGetReward {
	player := c.player.GetPlayerBySession(s)
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

func (c *MailRoute) ReqRead(s *network.Session, index int32, msg *protos.ReqMailRead) *protos.ResMailRead {
	player := c.player.GetPlayerBySession(s)
	code := c.service.Read(player, msg.Id)
	return &protos.ResMailRead{
		Code: int32(code),
	}
}
