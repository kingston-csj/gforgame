package player

import (
	"io/github/gforgame/domain"
	"io/github/gforgame/examples/constants"
	"time"
)

type Mailbox struct {
	Mails map[string]*Mail `json:"mails"`
	// 服务器邮件最大id
	ServerMailMaxId string `json:"serverMailMaxId"`
}

func (m *Mailbox) AfterLoad() {
	if m.Mails == nil {
		m.Mails = make(map[string]*Mail)
	}
}

func (m *Mailbox) AddSevMail(serverMail *ServerMail) {
	mail := &Mail{
		Id:      serverMail.Id,
		Title:   serverMail.Title,
		Content: serverMail.Content,
        Rewards: []domain.RewardDefLite{},
        Status:  constants.MailStatusUnread,
        Time:    time.Now().UnixMilli(),
    }
	m.AddMail(mail)
}

func (m *Mailbox) AddMail(mail *Mail) {
	if len(m.Mails) >= constants.MAIL_MAX_CAPACITY {
		//从早到晚，删除第一封已读已领取/已过期
		firstMail := m.GetMailList()[0]
		var toRemove *Mail
		for _, mail := range m.Mails {
			// 如果邮件已读已领取/已过期，则删除
			if mail.Status == constants.MailStatusReceived || mail.IsExpired() {
				toRemove = mail
				break
			}
			// 无奖励，且已读
			if len(mail.Rewards) == 0 && mail.Status == constants.MailStatusRead {
				toRemove = mail
				break
			}
		}
		if toRemove == nil {
			toRemove = firstMail
		}
		delete(m.Mails, toRemove.Id)
	}
	m.Mails[mail.Id] = mail
}

func (m *Mailbox) GetMail(id string) *Mail {
	return m.Mails[id]
}

func (m *Mailbox) GetMailList() []*Mail {
	mailList := make([]*Mail, 0, len(m.Mails))
	for _, mail := range m.Mails {
		mailList = append(mailList, mail)
	}
	return mailList
}

func (m *Mailbox) MarkReceivedServerMail(id string) {
	m.ServerMailMaxId = id
}

func (m *Mailbox) HasReceivedServerMail(id string) bool {
	return m.ServerMailMaxId >= id
}

func (m *Mailbox) DeleteMail(id string) {
	delete(m.Mails, id)
}
