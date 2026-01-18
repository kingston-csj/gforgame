package player

import (
	"io/github/gforgame/examples/constants"
	"time"
)

type RechargeBox struct {
	// 七日活动：是否已激活付费奖励
	ActivatedQiRiPay int32 `json:"activatedQiRiPay"`
	// 充值累计次数
	RechargeTimes map[int32]int32 `json:"rechargeTimes"`
	// 通行证活动：是否已激活付费奖励
	ActivatedPassPay int32 `json:"activatedPassPay"`
	// 首充时间
	FirstRechargeTime int64 `json:"firstRechargeTime"`
	// 银月卡信息
	SilverCard MonthlyCardVo `json:"silverCard"`
	// 金月卡信息
	GoldCard MonthlyCardVo `json:"goldCard"`
}

func (r *RechargeBox) GetOrCreateMonthlyCardVo(monthCardType int32) *MonthlyCardVo {
	switch monthCardType {
	case constants.MonthCardTypeSilver:
		return &r.SilverCard
	case constants.MonthCardTypeGold:
		return &r.GoldCard
	default:
		return nil
	}
}

type MonthlyCardVo struct {
	ExpiredTime int64 `json:"expiredTime"` // 月卡过期时间
}

func (r *MonthlyCardVo) IsActivated() bool {
	return r.ExpiredTime > time.Now().Unix()
}