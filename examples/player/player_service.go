package player

import (
	"errors"
	"sync"

	"io/github/gforgame/examples/attribute"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrCast     = errors.New("cast exception")
)

var (
	instance *PlayerService
	once     sync.Once
)

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

func (ps *PlayerService) GetOrCreatePlayer(playerId string) *playerdomain.Player {
	player := ps.GetPlayer(playerId)
	if player == nil {
		player = &playerdomain.Player{}
		player.Id = playerId
		player.Name = ""
		player.Level = 1
		ps.SavePlayer(player)
	}
	return player
}

func (ps *PlayerService) SavePlayer(player *playerdomain.Player) {
	cache, _ := context.CacheManager.GetCache("player")
	cache.Set(player.GetID(), player)
	context.DbService.SaveToDb(player)
}

func (ps *PlayerService) NotifyPurseChange(player *playerdomain.Player) {
	resPurse := &protos.PushPurseInfo{}
	resPurse.Diamond = player.Purse.Diamond
	resPurse.Gold = player.Purse.Gold
	io.NotifyPlayer(player, resPurse)
}

func (ps *PlayerService) RecCalculatePlayerAttr(player *playerdomain.Player) {
	// 主公本身属性+所有上阵英雄属性
	attrContainer := attribute.NewAttrBox()

	for _, hero := range player.HeroBox.Heros {
		attrContainer.AddAttrs(hero.AttrBox.GetAttrs())
	}

}
