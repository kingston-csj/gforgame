package gm

import (
	commonerrors "github.com/forfun/gforgame/common/errors"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/system"
)

type SystemGmHandler struct{}

func NewSystemGmHandler() *SystemGmHandler {
	return &SystemGmHandler{}
}

func (h *SystemGmHandler) RegisterTo(gm *GmService) {
	gm.Register("help", "查看所有GM命令", "help", gm.handleHelp)
	gm.Register("daily_reset", "触发每日重置", "daily_reset", h.handleDailyReset)
}

func (h *SystemGmHandler) handleDailyReset(player *playerdomain.Player, params string) *commonerrors.BusinessError {
	system.PerformDailyUpdate()
	return nil
}
