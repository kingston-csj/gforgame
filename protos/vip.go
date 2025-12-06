package protos

// PushVipQueryInfo 推送vip查询信息
type PushVipQueryInfo struct {
	VipLevel int32 `json:"vip_level"`  //当前vip等级
	ExpiredTime int64 `json:"expired_time"` //过期时间
	RechargeRmb float32 `json:"recharge_rmb"` //本周期累计充值金额
	PeriodRechargeRmb float32 `json:"period_recharge_rmb"` //本周期充值金额
}
