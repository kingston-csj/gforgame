package route

import (
	"github.com/forfun/gforgame/internal/context"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/network"
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
