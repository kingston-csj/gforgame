package mail

import (
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/reward"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type MailController struct {
	network.Base
}

func NewMailController() *MailController {
	return &MailController{}
}

func (c *MailController) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		c.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (c *MailController) OnPlayerLogin(player *playerdomain.Player) {
	GetMailService().checkMailsOnLogin(player)
}

func (c *MailController) ReqGetAllRewards(s *network.Session, index int, msg *protos.ReqMailGetAllRewards) *protos.ResMailGetAllRewards {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	mailList := player.Mailbox.GetMailList()
	andReward := reward.NewAndReward()
	for _, mail := range mailList {
		if mail.Status != constants.MailStatusReceived {
			mail.Status = constants.MailStatusReceived
			mailRewards := reward.ParseRewards(mail.Rewards)
			andReward.AddReward(mailRewards)
		}
	}
	andReward = andReward.Merge()
	andReward.Reward(player)
	context.EventBus.Publish(events.PlayerEntityChange, player)
	return &protos.ResMailGetAllRewards{
		Code: 0,
	}
}

func (c *MailController) ReqDeleteAll(s *network.Session, index int, msg *protos.ReqMailDeleteAll) *protos.ResMailDeleteAll {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	removed := make([]int64, 0)
	for _, mail := range player.Mailbox.GetMailList() {
		if mail.Status == constants.MailStatusReceived {
			// 删除
			player.Mailbox.DeleteMail(mail.Id)
			removed = append(removed, mail.Id)
		}
	}
	context.EventBus.Publish(events.PlayerEntityChange, player)

	return &protos.ResMailDeleteAll{
		Removed: removed,
	}
}

func (c *MailController) ReqGetReward(s *network.Session, index int, msg *protos.ReqMailGetReward) *protos.ResMailGetReward {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)

	mail := player.Mailbox.GetMail(msg.Id)
	if mail == nil {
		return &protos.ResMailGetReward{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}
	mailRewards := reward.ParseRewards(mail.Rewards)
	mail.Status = constants.MailStatusReceived
	mailRewards.Reward(player)

	context.EventBus.Publish(events.PlayerEntityChange, player)

	return &protos.ResMailGetReward{
		Code: 0,
	}
}

func (c *MailController) ReqRead(s *network.Session, index int, msg *protos.ReqMailRead) *protos.ResMailRead {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	mail := player.Mailbox.GetMail(msg.Id)
	if mail == nil {
		return &protos.ResMailRead{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}
	mail.Status = constants.MailStatusRead

	context.EventBus.Publish(events.PlayerEntityChange, player)

	return &protos.ResMailRead{}
}
