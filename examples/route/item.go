package route

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type ItemRoute struct {
	network.Base
}

func NewItemRoute() *ItemRoute {
	return &ItemRoute{}
}

func (ps *ItemRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		ps.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *ItemRoute) OnPlayerLogin(player *playerdomain.Player) {
	// 发送背包信息
	resBackpack := &protos.PushBackpackInfo{
		Items: []protos.ItemInfo{},
	}
	if player.Backpack != nil {
		for _, item := range player.Backpack.Items {
			resBackpack.Items = append(resBackpack.Items, item.ToVo())
		}
	}
	io.NotifyPlayer(player, resBackpack)

	// 发送货币信息
	player.NotifyPurseChange()
}
