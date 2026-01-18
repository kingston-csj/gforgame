package mail

import (
	"time"

	"io/github/gforgame/domain"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/reward"

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

func (s *MailService) CheckMailsOnLogin(player *playerdomain.Player) {
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

func (s *MailService) Read(player *playerdomain.Player, mailId int64) int {
	mail := player.Mailbox.GetMail(mailId)
	if mail == nil {
		return constants.I18N_COMMON_NOT_FOUND
	}
	mail.Status = constants.MailStatusRead

	context.EventBus.Publish(events.PlayerEntityChange, player)
	return 0
}

func (s *MailService) TakeReward(player *playerdomain.Player, mailId int64) (int, []*protos.RewardVo) {
	mail := player.Mailbox.GetMail(mailId)
	if mail == nil {
		return constants.I18N_COMMON_NOT_FOUND, nil
	}
	mailRewards := reward.ParseRewards(mail.Rewards)
	mail.Status = constants.MailStatusReceived
	mailRewards.Reward(player, constants.ActionType_MailGetReward)

	context.EventBus.Publish(events.PlayerEntityChange, player)
	return 0, reward.ToRewardVos(mailRewards)
}

func (s *MailService) TakeAllRewards(player *playerdomain.Player) []*protos.RewardVo {
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
	andReward.Reward(player, constants.ActionType_MailGetAll)
	context.EventBus.Publish(events.PlayerEntityChange, player)
	return reward.ToRewardVos(andReward)
}

func (s *MailService) DeleteAll(player *playerdomain.Player) []int64 {
	removed := make([]int64, 0)
	for _, mail := range player.Mailbox.GetMailList() {
		if mail.Status == constants.MailStatusReceived {
			// 删除
			player.Mailbox.DeleteMail(mail.Id)
			removed = append(removed, mail.Id)
		}
	}
	context.EventBus.Publish(events.PlayerEntityChange, player)
	return removed
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
            Rewards: []domain.RewardDefLite{
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
            Rewards: []domain.RewardDefLite{
                {
                    Type:  "currency",
                    Value: "gold=100",
                },
            },
        })

	}

	for _, mail := range player.Mailbox.GetMailList() {
		rewardVo := make([]protos.RewardVo, 0, len(mail.Rewards))
		for _, reward := range mail.Rewards {
			rewardVo = append(rewardVo, protos.RewardVo{
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
