package events

const (
	PlayerLogin        = "player_login"
	PlayerAfterLoad    = "player_after_load"
	PlayerEntityChange = "player_entity_change"
	PlayerAttrChange   = "player_attr_change"
	// 客户端加载完成
	PlayerLoadingFinish = "player_loading_finish"
	PlayerDailyReset    = "player_daily_reset"

	ItemConsume = "item_consume"

	HeroGain = "hero_gain"

	SystemDailyReset   = "system_daily_reset"
	SystemWeeklyReset  = "system_weekly_reset"
	SystemMonthlyReset = "system_monthly_reset"

	// 招募
	Recruit = "player_recruit"
)
