package config

type BuffData struct {
	Id int32 `json:"id" excel:"id"`
	// 持续回合数
	Duration int32 `json:"duration" excel:"duration"`
	// 叠加关系
	Relation int32 `json:"relation" excel:"relation"`
	// 最大叠加层数
	Layer int32 `json:"layer" excel:"layer"`
	// 类型 1为属性，2为状态
	Type int32 `json:"type" excel:"type"`
	// 参数
	Params string `json:"params" excel:"params"`
}
