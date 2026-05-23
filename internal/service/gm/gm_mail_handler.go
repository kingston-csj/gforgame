package gm

import (
	commonerrors "github.com/forfun/gforgame/common/errors"
	"github.com/forfun/gforgame/common/util/conv"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/service/mail"
)

type MailGmHandler struct {
	mail *mail.MailService
}

func NewMailGmHandler(mailService *mail.MailService) *MailGmHandler {
	return &MailGmHandler{
		mail: mailService,
	}
}

func (h *MailGmHandler) RegisterTo(gm *GmService) {
	gm.Register("mail", "添加邮件", "mail 1001", h.handleAddMail)
}

func (h *MailGmHandler) handleAddMail(player *playerdomain.Player, params string) *commonerrors.BusinessError {
	h.mail.SendSimpleMail(player, conv.Int32Value(params))
	return nil
}
