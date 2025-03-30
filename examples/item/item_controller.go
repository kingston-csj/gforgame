package item

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/utils"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type ItemController struct {
	network.Base
}

func NewItemController() *ItemController {
	return &ItemController{}
}

func (ps *ItemController) Init() {
	network.RegisterMessage(protos.CmdItemResBackpackInfo, &protos.ResBackpackInfo{})

	context.EventBus.Subscribe("player_login", func(data interface{}) {
		ps.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *ItemController) OnPlayerLogin(player *playerdomain.Player) {
	resBackpack := &protos.ResBackpackInfo{}
	if player.Backpack != nil {
		// 临时处理，后续采用事件驱动
		for id, count := range player.Backpack.Items {
			resBackpack.Items = append(resBackpack.Items, protos.ItemInfo{
				Id:    int32(id),
				Count: int32(count),
			})
		}
	}
	utils.NotifyPlayer(player, "item_backpack_info", resBackpack)
}
