package gm

import (
	commonerrors "github.com/forfun/gforgame/common/errors"
	"github.com/forfun/gforgame/common/util/conv"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/service/recharge"
)

type RechargeGmHandler struct {
	recharge *recharge.RechargeService
}

func NewRechargeGmHandler(rechargeService *recharge.RechargeService) *RechargeGmHandler {
	return &RechargeGmHandler{
		recharge: rechargeService,
	}
}

func (h *RechargeGmHandler) RegisterTo(gm *GmService) {
	gm.Register("recharge", "模拟充值", "recharge 1", h.handleRecharge)
}

func (h *RechargeGmHandler) handleRecharge(player *playerdomain.Player, params string) *commonerrors.BusinessError {
	rechargeId, _ := conv.StringToInt32(params)
	h.recharge.Recharge(player, rechargeId)
	return nil
}
