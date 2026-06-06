package route

import (
	"github.com/forfun/gforgame/internal/service/item"
	"github.com/forfun/gforgame/internal/service/player"
)

type ItemRoute struct {
	service *item.ItemService
	player  *player.PlayerService
}

func NewItemRoute(service *item.ItemService, playerService *player.PlayerService) *ItemRoute {
	return &ItemRoute{
		service: service,
		player:  playerService,
	}
}



