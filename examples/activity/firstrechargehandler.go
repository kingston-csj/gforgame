package activity

import (
	"io/github/gforgame/examples/domain/player"
)

type FirstRechargeActivityHandler struct {
	*BaseActivityHandler
}

func (d *FirstRechargeActivityHandler) GetBaseHandler() *BaseActivityHandler {
	return d.BaseActivityHandler
}

func (d *FirstRechargeActivityHandler) LoadActivityInfo(player *player.Player) {
	 
}

func NewFirstRechargeActivityHandler(sched *ActivityScheduler) *FirstRechargeActivityHandler {
	baseHandler := &BaseActivityHandler{
		ActivitySched: sched,
		ActivityId:    1001,
		OnActivityStart: func() error {
			// 首次充值活动的专属启动逻辑
			return nil
		},
		OnActivityEnd: func() error {
			// 首次充值活动的专属结束逻辑
			return nil
		},
	}

	return &FirstRechargeActivityHandler{
		BaseActivityHandler: baseHandler,
	}
}