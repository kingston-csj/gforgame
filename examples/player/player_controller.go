package player

import (
	"fmt"

	"io/github/gforgame/common"
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/consume"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/hero"
	"io/github/gforgame/examples/system"

	"io/github/gforgame/logger"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"

	"gorm.io/gorm"
)

type PlayerController struct {
	network.Base
}

func NewPlayerController() *PlayerController {
	return &PlayerController{}
}

func (ps *PlayerController) Init() {
	network.RegisterMessage(protos.CmdPlayerReqLogin, &protos.ReqPlayerLogin{})
	network.RegisterMessage(protos.CmdPlayerResLogin, &protos.ResPlayerLogin{})

	network.RegisterMessage(protos.CmdPlayerReqCreate, &protos.ReqPlayerCreate{})
	network.RegisterMessage(protos.CmdPlayerResCreate, &protos.ResPlayerCreate{})
	network.RegisterMessage(protos.CmdPlayerReqLoadingFinish, &protos.ReqPlayerLoadingFinish{})
	network.RegisterMessage(protos.CmdPlayerReqUpLevel, &protos.ReqPlayerUpLevel{})
	network.RegisterMessage(protos.CmdPlayerResUpLevel, &protos.ResPlayerUpLevel{})
	network.RegisterMessage(protos.CmdPlayerReqUpStage, &protos.ReqPlayerUpStage{})
	network.RegisterMessage(protos.CmdPlayerResUpStage, &protos.ResPlayerUpStage{})
	network.RegisterMessage(protos.CmdPlayerPushFightChange, &protos.PushPlayerFightChange{})

	// 自动建表
	err := mysqldb.Db.AutoMigrate(&playerdomain.Player{})
	if err != nil {
		panic(err)
	}

	// 缓存数据读取
	dbLoader := func(key string) (interface{}, error) {
		var p playerdomain.Player
		result := mysqldb.Db.First(&p, "id=?", key)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// 未找到记录
				return nil, nil
			}
		}
		p.AfterFind(nil)
		context.EventBus.Publish(events.PlayerAfterLoad, &p)
		return &p, nil
	}
	context.CacheManager.Register("player", dbLoader)

	context.EventBus.Subscribe(events.PlayerEntityChange, func(data interface{}) {
		GetPlayerService().SavePlayer(data.(*playerdomain.Player))
	})

	context.EventBus.Subscribe(events.PlayerAttrChange, func(data interface{}) {
		GetPlayerService().refreshFighting(data.(*playerdomain.Player))
	})

	// 在线玩家每日重置
	context.EventBus.Subscribe(events.SystemDailyReset, func(data interface{}) {
		allSessions := network.GetAllOnlinePlayerSessions()
		for _, s := range allSessions {
			s.AsynTasks <- func() {
				player := network.GetPlayerBySession(s).(*playerdomain.Player)
				player.DailyReset.Reset(data.(int64))
				GetPlayerService().SavePlayer(player)
			}
		}
	})
}

func (ps *PlayerController) ReqLogin(s *network.Session, index int, msg *protos.ReqPlayerLogin) {
	if util.IsBlankString(msg.PlayerId) {
		s.Send(&protos.ResPlayerLogin{Code: constants.I18N_COMMON_ILLEGAL_PARAMS}, index)
		return
	}
	// 是否是新角色
	newCreated := GetPlayerService().GetPlayerProfileById(msg.PlayerId) == nil
	player := GetPlayerService().GetOrCreatePlayer(msg.PlayerId)
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

	// 离线，登录触发每日重置检测
	dailyReset := system.GetDailyReset().GetValue().(int64)
	if player.DailyReset.LastDailyReset > 0 && player.DailyReset.LastDailyReset < dailyReset {
		player.DailyReset.Reset(dailyReset)
		GetPlayerService().SavePlayer(player)
	}
	// 添加session
	network.AddSession(s, player)

	// 客户端红点系统，要求服务器先下发所有基础数据后，客户端再切到主界面
	context.EventBus.Publish(events.PlayerLogin, player)

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

func (ps *PlayerController) ReqLoadingFinish(s *network.Session, index int, msg *protos.ReqPlayerLoadingFinish) {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	context.EventBus.Publish(events.PlayerLoadingFinish, player)
}

func (ps *PlayerController) ReqCreate(s *network.Session, msg *protos.ReqPlayerCreate) {
	id := util.GetNextID()
	player := &playerdomain.Player{}
	player.Id = id
	player.Name = msg.Name
	mysqldb.Db.Create(&player)

	logger.Log(logger.Player, "Id", player.Id, "name", player.Name)
	fmt.Printf(player.Name)
}

func (ps *PlayerController) ReqPlayerUpLevel(s *network.Session, index int, msg *protos.ReqPlayerUpLevel) *protos.ResPlayerUpLevel {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	toLevel := msg.ToLevel
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
	consume.Consume(p)

	p.Level = toLevel
	GetPlayerService().refreshFighting(p)
	GetPlayerService().SavePlayer(p)

	return &protos.ResPlayerUpLevel{
		Code: 0,
	}
}

func (ps *PlayerController) ReqHeroUpStage(s *network.Session, index int, msg *protos.ReqPlayerUpStage) *protos.ResHeroUpStage {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)

	stageData := config.GetSpecificContainer[container.HeroStageContainer]("herostage").GetRecordByStage(p.Stage)
	if stageData == nil {
		return &protos.ResHeroUpStage{
			Code: constants.I18N_HERO_TIP4,
		}
	}

	if p.Level < stageData.MaxLevel {
		return &protos.ResHeroUpStage{
			Code: constants.I18N_HERO_TIP3,
		}
	}

	stageData = config.GetSpecificContainer[container.HeroStageContainer]("herostage").GetRecordByStage(p.Stage + 1)
	if stageData == nil {
		return &protos.ResHeroUpStage{
			Code: constants.I18N_HERO_TIP4,
		}
	}

	costItem := consume.ItemConsume{
		ItemId: constants.GAME_UPSTAGE_ITEM_ID,
		Amount: stageData.Cost,
	}
	err := costItem.Verify(p)
	if err != nil {
		return &protos.ResHeroUpStage{
			Code: int32(err.(*common.BusinessRequestException).Code()),
		}
	}
	costItem.Consume(p)

	p.Stage = p.Stage + 1

	GetPlayerService().refreshFighting(p)
	GetPlayerService().SavePlayer(p)

	return &protos.ResHeroUpStage{
		Code: 0,
	}
}
