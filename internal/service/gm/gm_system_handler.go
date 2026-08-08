package gm

import (
	commonerrors "github.com/forfun/gforgame/common/errors"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	playerservice "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/internal/system"
)

type SystemGmHandler struct {
	system        *system.SystemService
	playerService *playerservice.PlayerService
}

func NewSystemGmHandler(systemService *system.SystemService, playerService *playerservice.PlayerService) *SystemGmHandler {
	return &SystemGmHandler{system: systemService, playerService: playerService}
}

func (h *SystemGmHandler) RegisterTo(gm *GmService) {
	gm.Register("help", "查看所有GM命令", "help", gm.handleHelp)
	gm.Register("daily_reset", "触发每日重置", "daily_reset", h.handleDailyReset)
}

func (h *SystemGmHandler) handleDailyReset(player *playerdomain.Player, params string) *commonerrors.BusinessError {
	resetTime := h.system.GetDailyReset().ResetTime
	h.playerService.DailyReset(player, resetTime)
	return nil
}
