package player

import "io/github/gforgame/examples/attribute"

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
}
