package mail

import (
	"time"

	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/domain/config/item"

	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
)

var (
	instance    *MailService = &MailService{}
	serverMails              = make(map[int64]*playerdomain.ServerMail)
)

type MailService struct{}

func GetMailService() *MailService {
	return instance
}

func (s *MailService) AddServerMail(serverMail *playerdomain.ServerMail) {
	serverMails[serverMail.Id] = serverMail
}

func (s *MailService) checkMailsOnLogin(player *playerdomain.Player) {
	s.checkServerMails(player)
	s.notifyMails(player)
}

func (s *MailService) checkServerMails(player *playerdomain.Player) {
	mailbox := player.Mailbox
	for _, serverMail := range serverMails {
		if !player.Mailbox.HasReceivedServerMail(serverMail.Id) {
			player.Mailbox.MarkReceivedServerMail(serverMail.Id)
		}
		mailbox.AddSevMail(serverMail)
	}
}

func (s *MailService) notifyMails(player *playerdomain.Player) {
	mailList := make([]protos.MailVo, 0, 10)

	if len(player.Mailbox.GetMailList()) == 0 {
		// 加入一些测试数据
		player.Mailbox.AddMail(&playerdomain.Mail{
			Id:      1,
			Title:   "测试邮件1",
			Content: "测试邮件1内容",
			Time:    time.Now().UnixMilli(),
			Status:  constants.MailStatusUnread,
			Rewards: []item.RewardDef{
				{
					Type:  "item",
					Value: "2003=1",
				},
			},
		})

		player.Mailbox.AddMail(&playerdomain.Mail{
			Id:      2,
			Title:   "测试邮件2",
			Content: "测试邮件2内容",
			Time:    time.Now().UnixMilli(),
			Status:  constants.MailStatusUnread,
			Rewards: []item.RewardDef{
				{
					Type:  "currency",
					Value: "gold=100",
				},
			},
		})

	}

	for _, mail := range player.Mailbox.GetMailList() {
		rewardVo := make([]protos.RewardInfo, 0, len(mail.Rewards))
		for _, reward := range mail.Rewards {
			rewardVo = append(rewardVo, protos.RewardInfo{
				Type:  reward.Type,
				Value: reward.Value,
			})
		}
		mailVo := protos.MailVo{
			Id:      mail.Id,
			Title:   mail.Title,
			Content: mail.Content,
			Time:    mail.Time,
			Status:  mail.Status,
			Rewards: rewardVo,
		}
		mailList = append(mailList, mailVo)
	}
	resMailList := &protos.PushMailAll{
		Mails: mailList,
	}

	io.NotifyPlayer(player, resMailList)
}
