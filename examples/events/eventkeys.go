package events

const (
	PlayerLogin        = "player_login"
	PlayerAfterLoad    = "player_after_load"
	PlayerEntityChange = "player_entity_change"
	PlayerAttrChange   = "player_attr_change"
	// 客户端加载完成
	PlayerLoadingFinish = "player_loading_finish"
	PlayerDailyReset    = "player_daily_reset"

	Recharge = "recharge"

	ItemConsume = "item_consume"
	MallBuy     = "mall_buy"

	HeroGain    = "hero_gain"
	HeroLevelUp = "hero_level_up"

	SystemDailyReset   = "system_daily_reset"
	SystemWeeklyReset  = "system_weekly_reset"
	SystemMonthlyReset = "system_monthly_reset"

	// 招募
	Recruit = "player_recruit"
	// 客户端事件
	ClientDiyEvent = "client_diy_event"

	HeroEntrust = "hero_entrust"

	// 竞技场积分改变
	AreaScoreChanged = "area_score_changed"
	// 竞技场通过
	PassArena = "pass_arena"
)
