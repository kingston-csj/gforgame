package route

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/service/mail"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type MailRoute struct {
	network.Base
	service *mail.MailService
}

func NewMailRoute() *MailRoute {
	return &MailRoute{
		service: mail.GetMailService(),
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
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	rewardVos := c.service.TakeAllRewards(player)
	return &protos.ResMailGetAllRewards{
		Code: 0,
		Rewards: rewardVos,
	}
}

func (c *MailRoute) ReqDeleteAll(s *network.Session, index int32, msg *protos.ReqMailDeleteAll) *protos.ResMailDeleteAll {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	removed := c.service.DeleteAll(player)
	return &protos.ResMailDeleteAll{
		Removed: removed,
	}
}

func (c *MailRoute) ReqGetReward(s *network.Session, index int32, msg *protos.ReqMailGetReward) *protos.ResMailGetReward {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
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
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	code := c.service.Read(player, msg.Id)
	return &protos.ResMailRead{
		Code: int32(code),
	}
}
