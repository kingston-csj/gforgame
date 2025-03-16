package player

import (
	"errors"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/network"
	"sync"
)

var ErrNotFound = errors.New("record not found")
var ErrCast = errors.New("cast exception")

var instance *PlayerService
var once sync.Once

type PlayerService struct {
	network.Base
}

func GetPlayerService() *PlayerService {
	once.Do(func() {
		instance = &PlayerService{}
	})
	return instance
}

func (ps *PlayerService) GetPlayer(playerId string) *playerdomain.Player {
	cache, _ := context.CacheManager.GetCache("player")
	cacheEntity, err := cache.Get(playerId)
	if err != nil {
		return nil
	}
	if cacheEntity == nil {
		return nil
	}
	player, _ := cacheEntity.(*playerdomain.Player)
	return player
}

func (ps *PlayerService) SavePlayer(player *playerdomain.Player) {
	cache, _ := context.CacheManager.GetCache("player")
	cache.Set(player.GetId(), player)
	context.DbService.SaveToDb(player)
}
