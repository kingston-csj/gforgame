package player

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"io/github/gforgame/common"
	"io/github/gforgame/db"
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/camp"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/consume"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/fight/attribute"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/service/hero"
	"io/github/gforgame/examples/system"
	"io/github/gforgame/logger"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"
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

func (ps *PlayerService) DoLogin(playerId string, s *network.Session, index int32) {
	// 是否是新角色
	newCreated := ps.GetPlayerProfileById(playerId) == nil
	player := ps.GetOrCreatePlayer(playerId)
	if !newCreated {
		oldSession := network.GetSessionByPlayerId(player.Id)
		if oldSession != nil {
			if oldSession == s {
				logger.Info("玩家重复登录[" + player.Id + "]")
			} else {
				// 旧会话存在，关闭旧会话
				logger.Info("玩家顶号登录[" + player.Id + "]")
				oldSession.SendAndClose(&protos.PushReplacingLogin{})
			}
		}
	}
	fmt.Println("登录成功，id为：", player.Id)

	// 添加session
	network.AddSession(s, player)
	
	// 离线，登录触发每日重置检测
	dailyReset := system.GetDailyReset().GetValue().(int64)
	if player.DailyReset.LastDailyReset > 0 && player.DailyReset.LastDailyReset < dailyReset {
		ps.dailyReset(player, dailyReset)
	} else {
		ps.PushDailyResetInfo(player)
	}
	

	// 客户端红点系统，要求服务器先下发所有基础数据
	// 异步推送
	go func(){
		context.EventBus.Publish(events.PlayerLogin, player)
		// 客户端再切到主界面
		s.SendWithoutIndex(&protos.PushLoadComplete{})
	}()

	s.Send(&protos.ResPlayerLogin{
		Code:       0,
		PlayerId:   player.Id,
		NewCreate:  newCreated,
		CreateTime: player.CreateTime,
		Head:       player.Head,
		Level:      player.Level,
		Name:       player.Name,
		Fighting:   player.Fight,
		Camp:       player.Camp,
	}, index)
}

func (ps *PlayerService) Create(name string, camp int32) *playerdomain.Player {
	id := util.GetNextID()
	player := &playerdomain.Player{}
	player.Id = id
	player.Name = name
	player.Camp = camp
	mysqldb.Db.Create(&player)

	logger.Log(logger.Player, "Id", player.Id, "name", player.Name)
	fmt.Printf(player.Name)

	return player
}

func (ps *PlayerService) DoUpLevel(p *playerdomain.Player, toLevel int32)  *protos.ResPlayerUpLevel {
	stageData := config.GetSpecificContainer[container.HeroStageContainer]("herostage").GetRecordByStage(p.Stage)
	if stageData == nil {
		return &protos.ResPlayerUpLevel{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}
	if p.Level >= stageData.MaxLevel {
		return &protos.ResPlayerUpLevel{
			Code: constants.I18N_HERO_TIP2,
		}
	}

	if toLevel <= p.Level {
		return &protos.ResPlayerUpLevel{
			Code: constants.I18N_COMMON_ILLEGAL_PARAMS,
		}
	}

	costGold := hero.GetHeroService().CalcTotalUpLevelConsume(p.Level, toLevel)
	if !p.Purse.IsEnoughGold(costGold) {
		return &protos.ResPlayerUpLevel{
			Code: constants.I18N_GOLD_NOT_ENOUGH,
		}
	}

	consume := consume.CurrencyConsume{
		Currency: "gold",
		Amount:   costGold,
	}
	err := consume.Verify(p)
	if err != nil {
		return &protos.ResPlayerUpLevel{
			Code: int32(err.(*common.BusinessRequestException).Code()),
		}	
	}
	consume.Consume(p, constants.ActionType_HeroUpLevel)

	p.Level = toLevel
	ps.RefreshFighting(p)
	ps.SavePlayer(p)

	return &protos.ResPlayerUpLevel{
		Code: 0,
	}
}

func (ps *PlayerService) DoUpStage(p *playerdomain.Player) *protos.ResPlayerUpStage {
	stageData := config.GetSpecificContainer[container.HeroStageContainer]("herostage").GetRecordByStage(p.Stage)
	if stageData == nil {
		return &protos.ResPlayerUpStage{
			Code: constants.I18N_HERO_TIP4,
		}
	}
	
	if p.Level < stageData.MaxLevel {
		return &protos.ResPlayerUpStage{
			Code: constants.I18N_HERO_TIP3,
		}
	}

	stageData = config.GetSpecificContainer[container.HeroStageContainer]("herostage").GetRecordByStage(p.Stage + 1)
	if stageData == nil {
		return &protos.ResPlayerUpStage{
			Code: constants.I18N_HERO_TIP4,
		}
	}

	costItem := consume.ItemConsume{
		ItemId: constants.GAME_UPSTAGE_ITEM_ID,
		Amount: stageData.Cost,
	}
	err := costItem.Verify(p)
	if err != nil {
		return &protos.ResPlayerUpStage{
			Code: int32(err.(*common.BusinessRequestException).Code()),
		}
	}
	costItem.Consume(p, constants.ActionType_HeroUpStage)

	p.Stage = p.Stage + 1

	ps.RefreshFighting(p)
	ps.SavePlayer(p)

	return &protos.ResPlayerUpStage{
		Code: 0,
	}
}

func (ps *PlayerService) RefreshFighting(player *playerdomain.Player) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Errorf("panic recovered: %v", r))
		}
	}()
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

 func (ps *PlayerService) dailyReset(player *playerdomain.Player, resetTime int64) {
	box := &playerdomain.DailyReset{
		LastDailyReset: resetTime,
	}
	player.DailyReset = box
	context.EventBus.Publish(events.PlayerDailyReset, player)
	ps.SavePlayer(player)
	ps.PushDailyResetInfo(player)
 }

  func (ps *PlayerService) PushDailyResetInfo(player *playerdomain.Player) {
	io.NotifyPlayer(player, &protos.PushDailyResetInfo{
		NormalRecruitTimes: player.HeroBox.RecruitSum,
	})
  }

