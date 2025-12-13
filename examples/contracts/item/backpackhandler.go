package item

import (
	playerdomain "io/github/gforgame/examples/domain/player"
	protos "io/github/gforgame/protos"
)

type BackpackHandler interface {

	// 通过道具模型id进行扣除（后端接口）
	UseByModelId(p *playerdomain.Player, itemId int32, count int32) error 

	// 通过道具uid进行扣除（后端接口）
	UseByUid(p *playerdomain.Player, itemUid string, count int32) (error, []protos.RewardVo)

	AddByModelId(p *playerdomain.Player, itemId int32, count int32) error;
}