package player

import (
	"errors"
	"sync"

	"io/github/gforgame/examples/attribute"
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/config"
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

func (ps *PlayerService) refreshFighting(player *playerdomain.Player) {
	ps.recomputeAttribute(player)
	fighting := 0
	for _, hero := range player.HeroBox.Heros {
		fighting += int(hero.Fight)
	}
	player.Fight = int32(fighting)
	io.NotifyPlayer(player, &protos.PushPlayerFightChange{
		Fight: player.Fight,
	})
}

func (ps *PlayerService) recomputeAttribute(player *playerdomain.Player) {
	attrContainer := attribute.NewAttrBox()
	// 主公等级属性
	heroLevelDataRecord := context.GetDataManager().GetRecord("herolevel", int64(player.Level))
	heroLevelData := heroLevelDataRecord.(config.HeroLevelData)
	attrContainer.AddAttrs(heroLevelData.GetHeroLevelAttrs())

	// 主公突破属性
	heroStageDataRecord := context.GetDataManager().GetRecord("herostage", int64(player.Stage))
	if heroStageDataRecord != nil {
		heroStageData := heroStageDataRecord.(config.HeroStageData)
		attrContainer.AddAttrs(heroStageData.Attrs)
	}

	player.AttrBox = attrContainer
}

func (ps *PlayerService) getHeroIdByCamp(camp int32) int32 {
	if camp == 1001 {
		return 1001
	}
	if camp == 1002 {
		return 1002
	}
	if camp == 1003 {
		return 1003
	}
	return 1004
}
