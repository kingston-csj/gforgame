package route

import (
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/service/item"
)

type ItemRoute struct {
	service    *item.ItemService
	playerrepo *playerrepo.PlayerRepository
}

func NewItemRoute(service *item.ItemService, playerRepo *playerrepo.PlayerRepository) *ItemRoute {
	return &ItemRoute{
		service:    service,
		playerrepo: playerRepo,
	}
}
