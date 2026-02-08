package player

import "io/github/gforgame/examples/fight/attribute"

type Hero struct {
	// 模型ID
	ModelId int32
	// 等级
	Level int32
	// 阶段
	Stage int32
	// 属性
	AttrBox *attribute.AttrBox `json:"-"` // 该字段不参与序列化
	// 战斗力
	Fight int32
	// 位置 1-5 0为空闲位置
	Position int32
	// 委托任务ID
	EntrustQuestId int32
}
