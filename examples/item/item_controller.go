package item

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/io"
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
	network.RegisterMessage(protos.CmdItemResPurseInfo, &protos.PushPurseInfo{})
	network.RegisterMessage(protos.CmdItemPushChanged, &protos.PushItemChanged{})

	context.EventBus.Subscribe(events.PlayerLoadingFinish, func(data interface{}) {
		ps.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *ItemController) OnPlayerLogin(player *playerdomain.Player) {
	// 发送背包信息
	resBackpack := &protos.ResBackpackInfo{}
	if player.Backpack != nil {
		for id, count := range player.Backpack.Items {
			resBackpack.Items = append(resBackpack.Items, protos.ItemInfo{
				Id:    int32(id),
				Count: int32(count),
			})
		}
	}
	io.NotifyPlayer(player, resBackpack)

	// 发送货币信息
	player.NotifyPurseChange()
}
