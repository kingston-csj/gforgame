package mail

import (
	"time"

	"io/github/gforgame/domain"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/reward"
	"io/github/gforgame/util"
	"io/github/gforgame/util/timeutil"

	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
)

var (
	instance    *MailService = &MailService{}
	serverMails              = make(map[string]*playerdomain.ServerMail)
)

// 邮件模块
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

func (s *MailService) Read(player *playerdomain.Player, mailId string) int {
	mail := player.Mailbox.GetMail(mailId)
	if mail == nil {
		return constants.I18N_COMMON_NOT_FOUND
	}
	mail.Status = constants.MailStatusRead

	context.EventBus.Publish(events.PlayerEntityChange, player)
	return 0
}

func (s *MailService) TakeReward(player *playerdomain.Player, mailId string) (int, []*protos.RewardVo) {
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

func (s *MailService) DeleteAll(player *playerdomain.Player) []string {
	removed := make([]string, 0)
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
func (s *MailService) SendSimpleMail(player *playerdomain.Player,id int32) {
	mailData := config.QueryById[configdomain.MailData](id)
	if mailData == nil {
		return
	}
	s.SendMail(player, id, "", "", nil, mailData.ValidTime)
}

func (s *MailService) SendSimpleMail2(player *playerdomain.Player,id int32, params ...string) {
	mailData := config.QueryById[configdomain.MailData](id)
	if mailData == nil {
		return
	}
	s.SendMail(player, id, "", "", nil, mailData.ValidTime, params...)
}
func (s *MailService) SendMail(player *playerdomain.Player,id int32, title string, content string, rewards []domain.RewardDefLite, validHours int32, params ...string) {
	mailId := util.GetNextID()
	mailData := config.QueryById[configdomain.MailData](id)
	if mailData != nil && validHours == 0{
		validHours = mailData.ValidTime
	}
	player.Mailbox.AddMail(&playerdomain.Mail{
		Id:      mailId,
		Title:   title,
		Content: content,
		Time:    time.Now().UnixMilli(),
		Status:  constants.MailStatusUnread,
		Rewards: rewards,
		Params: params,
		ExpiredTime: time.Now().UnixMilli() + int64(validHours) * timeutil.MILLIS_PER_SECOND,
	})
}

func (s *MailService) notifyMails(player *playerdomain.Player) {
	mailList := make([]protos.MailVo, 0, 10)

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
