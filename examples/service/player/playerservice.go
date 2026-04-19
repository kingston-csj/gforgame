package player

import (
	"errors"
	"fmt"
	"io/github/gforgame/common"
	"io/github/gforgame/common/container/hashmap"
	"io/github/gforgame/common/logger"
	"io/github/gforgame/common/trie"
	"io/github/gforgame/common/util"
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
	mysqldb "io/github/gforgame/examples/infra/persistence"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/protos"
	"io/github/gforgame/examples/system"
	"io/github/gforgame/network"
	"strings"
	"sync"
)

var (
	ErrCast  = errors.New("cast exception")
	instance *PlayerService
	once     sync.Once
)

// 玩家模块
type PlayerService struct {
	network.Base
	playerProfiles map[string]*playerdomain.PlayerProfile
	// 双向map, id -> name
	idNameMapper *hashmap.SyncDualHashMap[string, string]
	// 玩家名称字典树
	nameDict *trie.TrieDictionary
}

func GetPlayerService() *PlayerService {
	once.Do(func() {
		instance = &PlayerService{
			playerProfiles: make(map[string]*playerdomain.PlayerProfile),
			idNameMapper:   hashmap.NewSyncDualHashMap[string, string](),
			nameDict:       trie.NewTrieDictionary(),
		}
	})
	return instance
}

// LoadPlayerProfile 加载玩家概况数据
func (ps *PlayerService) LoadPlayerProfile() {
	var profiles []*playerdomain.PlayerProfile
	err := mysqldb.Db.Model(&playerdomain.Player{}).Select("id, name, level, camp, fight").Scan(&profiles).Error
	if err != nil {
		panic(err)
	}

	for _, profile := range profiles {
		ps.playerProfiles[profile.Id] = profile
		ps.idNameMapper.Put(profile.Id, profile.Name)
	}
}

func (ps *PlayerService) GetPlayerProfileById(playerId string) *playerdomain.PlayerProfile {
	return ps.playerProfiles[playerId]
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

func (ps *PlayerService) GetPlayerBySession(session *network.Session) *playerdomain.Player {
	playerID, ok := network.GetPlayerIDBySession(session)
	if !ok {
		return nil
	}
	return ps.GetPlayer(playerID)
}

func (ps *PlayerService) GetPlayerByPlayerId(playerID string) *playerdomain.Player {
	return ps.GetPlayer(playerID)
}

func (ps *PlayerService) GetOrCreatePlayer(playerId string) *playerdomain.Player {
	player := ps.GetPlayer(playerId)
	if player == nil {
		player = &playerdomain.Player{}
		player.Id = playerId
		player.Camp = camp.Camp_Hao
		player.AfterFind(nil)
		initPlayer(player)
		ps.SavePlayer(player)
	}
	return player
}

func initPlayer(player *playerdomain.Player) {
	player.Name = instance.RandomName()
	player.Level = 1
	player.Stage = 1
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
	s.SetAttr("id", player.Id)

	// 添加session
	network.AddSession(s, player.Id)

	// 客户端红点系统，要求服务器先下发所有基础数据
	// 异步推送
	go func() {
		// 离线，登录触发每日重置检测
		dailyReset := system.GetDailyReset().GetValue().(int64)
		if player.DailyReset.LastDailyReset > 0 && player.DailyReset.LastDailyReset < dailyReset {
			ps.DailyReset(player, dailyReset)
		} else {
			ps.refreshDailyInfo(player)
		}
		weeklyReset := system.GetWeeklyReset().GetValue().(int64)
		if player.WeeklyReset.LastWeeklyReset > 0 && player.WeeklyReset.LastWeeklyReset < weeklyReset {
			ps.WeeklyReset(player, weeklyReset)
		} else {
			ps.refreshWeeklyInfo(player)
		}
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
		Stage:      player.Stage,
	}, index)
}

func (ps *PlayerService) Create(name string, camp int32) *playerdomain.Player {
	id := util.GetNextID()
	if util.IsEmptyString(name) {
		name = ps.RandomName()
	}
	player := &playerdomain.Player{}
	player.Id = id
	player.Name = name
	player.Camp = camp
	mysqldb.Db.Create(&player)

	logger.Log(constants.LoggerPlayer, "Id", player.Id, "name", player.Name)
	fmt.Printf(player.Name)

	return player
}

func (ps *PlayerService) DoUpLevel(p *playerdomain.Player, toLevel int32) *protos.ResPlayerUpLevel {
	stageData := config.GetSpecificContainer[*container.HeroStageContainer]().GetRecordByStage(p.Stage)
	if stageData == nil {
		return &protos.ResPlayerUpLevel{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}
	if p.Level >= stageData.MaxLevel {
		return &protos.ResPlayerUpLevel{
			Code: constants.I18N_COMMON_ILLEGAL_PARAMS,
		}
	}

	if toLevel <= p.Level {
		return &protos.ResPlayerUpLevel{
			Code: constants.I18N_COMMON_ILLEGAL_PARAMS,
		}
	}

	costGold := calcTotalUpLevelConsume(p.Level, toLevel)
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

func calcTotalUpLevelConsume(fromLevel int32, toLevel int32) int32 {
	levelContainer := config.GetSpecificContainer[*container.PlayerLevelContainer]()
	total := int32(0)
	for i := fromLevel; i < toLevel; i++ {
		levelData := levelContainer.GetLevelData(i, i)
		if levelData != nil {
			total += levelData.Cost
		}
	}
	return total
}

func (ps *PlayerService) DoUpStage(p *playerdomain.Player) *protos.ResPlayerUpStage {
	stageData := config.GetSpecificContainer[*container.HeroStageContainer]().GetRecordByStage(p.Stage)
	if stageData == nil {
		return &protos.ResPlayerUpStage{
			Code: constants.I18N_HERO_TIP4,
		}
	}

	if p.Level < stageData.MaxLevel {
		return &protos.ResPlayerUpStage{
			Code: constants.I18N_COMMON_ILLEGAL_PARAMS,
		}
	}

	stageData = config.GetSpecificContainer[*container.HeroStageContainer]().GetRecordByStage(p.Stage + 1)
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
			// logger.Error(fmt.Errorf("panic recovered: %v", r))
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
	stageContainer := config.GetSpecificContainer[*container.HeroStageContainer]()
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
	for _, profile := range ps.playerProfiles {
		if strings.Contains(profile.Name, name) {
			playerIds = append(playerIds, profile.Id)
		}
	}
	return playerIds
}

func (ps *PlayerService) DailyReset(player *playerdomain.Player, resetTime int64) {
	// 直接用新的数据替换，就不用为每个字段写重置逻辑了
	player.DailyReset = &playerdomain.DailyReset{
		LastDailyReset: resetTime,
	}
	player.DailyReset.AfterLoad()
	player.ExtendBox.AccumulatedLoginDays++
	context.EventBus.Publish(events.PlayerDailyReset, player)
	ps.SavePlayer(player)
	ps.refreshDailyInfo(player)
}

func (ps *PlayerService) refreshDailyInfo(player *playerdomain.Player) {
	push := &protos.PushDailyResetInfo{}
	io.NotifyPlayer(player, push)
}

func (ps *PlayerService) WeeklyReset(player *playerdomain.Player, resetTime int64) {
	// 直接用新的数据替换，就不用为每个字段写重置逻辑了
	player.WeeklyReset = &playerdomain.WeeklyReset{
		LastWeeklyReset: resetTime,
	}
	player.WeeklyReset.AfterLoad()
	context.EventBus.Publish(events.PlayerWeeklyReset, player)
	ps.SavePlayer(player)
	ps.refreshWeeklyInfo(player)
}

func (ps *PlayerService) refreshWeeklyInfo(player *playerdomain.Player) {
	push := &protos.PushWeeklyResetInfo{
		WeeklyRechargeSum: player.WeeklyReset.RechargeSum,
	}
	io.NotifyPlayer(player, push)
}
func (ps *PlayerService) RandomName() string {
	nameContainer := config.GetSpecificContainer[*container.NameContainer]()
	for i := 0; i < 10; i++ {
		name := nameContainer.GetRandomName()
		if _, ok := ps.idNameMapper.GetByKey(name); !ok {
			return name
		}
	}
	return util.GetNextID()
}

func (ps *PlayerService) EditPlayer(p *playerdomain.Player, head int32, name string) int32 {
	if p.Name != name {
		if _, ok := ps.idNameMapper.GetByKey(name); ok {
			return constants.I18N_PLAYER_NAME_REPEATED
		}
		oldName := p.Name
		p.Name = name
		ps.idNameMapper.Put(p.Id, name)
		ps.nameDict.DeleteNode(oldName)
		ps.nameDict.AddNode(name)
	}

	ps.playerProfiles[p.Id].Name = name
	p.Head = head
	ps.SavePlayer(p)
	return 0
}
