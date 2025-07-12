package player

import (
	"errors"
	"strings"
	"sync"

	"io/github/gforgame/db"
	"io/github/gforgame/examples/camp"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/fight/attribute"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

var (
	ErrNotFound    = errors.New("record not found")
	ErrCast        = errors.New("cast exception")
	instance       *PlayerService
	once           sync.Once
	playerProfiles map[string]*playerdomain.PlayerProfile = make(map[string]*playerdomain.PlayerProfile)
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

func (ps *PlayerService) LoadPlayerProfile() {
	var profiles []*playerdomain.PlayerProfile
	err := db.Db.Model(&playerdomain.Player{}).Select("id, name, level, camp, fight").Scan(&profiles).Error
	if err != nil {
		panic(err)
	}

	// 输出查询结果
	for _, profile := range profiles {
		playerProfiles[profile.Id] = profile
	}
}

func (ps *PlayerService) GetPlayerProfileById(playerId string) *playerdomain.PlayerProfile {
	return playerProfiles[playerId]
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
		player.Camp = camp.Camp_Hao
		player.AfterFind(nil)
		ps.SavePlayer(player)
	}
	return player
}

func (ps *PlayerService) SavePlayer(player *playerdomain.Player) {
	cache, _ := context.CacheManager.GetCache("player")
	cache.Set(player.GetId(), player)
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

	heroLevelData := config.QueryById[configdomain.HeroLevelData](player.Level)
	attrContainer.AddAttrs(heroLevelData.GetHeroLevelAttrs())

	// 主公突破属性
	stageContainer := config.QueryContainer[configdomain.HeroStageData, *container.HeroStageContainer]()
	stageData := stageContainer.GetRecordByStage(player.Stage)
	attrContainer.AddAttrs(stageData.GetHeroStageAttrs())

	player.AttrBox = attrContainer
}

func (ps *PlayerService) GetHeroIdByCamp(camp int32) int32 {
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

// 模糊搜索玩家(名字包含关键字)
func (ps *PlayerService) FuzzySearchPlayers(name string) []string {
	playerIds := make([]string, 0)
	for _, profile := range playerProfiles {
		if strings.Contains(profile.Name, name) {
			playerIds = append(playerIds, profile.Id)
		}
	}
	return playerIds
}
