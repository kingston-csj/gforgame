package config

type SkillData struct {
	Id   int32  `json:"id" excel:"id"`
	Name string `json:"name" excel:"name"`
	// 1主动2被动
	Type int32 `json:"type" excel:"type"`
	/// 技能效果
	EffectType int32  `json:"effect_type" excel:"effect_type"`
	Tips       string `json:"tips" excel:"tips"`
	// 技能目标选择器 1 自己 2 友方 3 敌方
	Selector int32 `json:"selector" excel:"selector"`
	// 技能范围，最多攻击几个目标
	AoeRange int32 `json:"aoe_range" excel:"aoe_range"`
	// 永久buff
	BuffId int32 `json:"buff_id" excel:"buff_id"`
	// 伤害系数（万分比）
	DamageRate int32 `json:"damage_rate" excel:"damage_rate"`
}
