package route

import (
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/service/player"

	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"

	"gorm.io/gorm"
)

type PlayerRoute struct {
	network.Base
	service *player.PlayerService
}

func NewPlayerRoute() *PlayerRoute {
	return &PlayerRoute{
		service: player.GetPlayerService(),
	}
}

func (ps *PlayerRoute) Init() {
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
		ps.service.SavePlayer(data.(*playerdomain.Player))
	})

	context.EventBus.Subscribe(events.PlayerAttrChange, func(data interface{}) {
		ps.service.RefreshFighting(data.(*playerdomain.Player))
	})

	// 在线玩家每日重置
	context.EventBus.Subscribe(events.SystemDailyReset, func(data interface{}) {
		allSessions := network.GetAllOnlinePlayerSessions()
		for _, s := range allSessions {
			s.AsynTasks <- func() {
				player := network.GetPlayerBySession(s).(*playerdomain.Player)
				ps.service.DailyReset(player, data.(int64))
			}
		}
	})
}

func (ps *PlayerRoute) ReqLogin(s *network.Session, index int32, msg *protos.ReqPlayerLogin) {
	if util.IsBlankString(msg.PlayerId) {
		s.Send(&protos.ResPlayerLogin{Code: constants.I18N_COMMON_ILLEGAL_PARAMS}, index)
		return
	}
	ps.service.DoLogin(msg.PlayerId, s, index)
}

func (ps *PlayerRoute) ReqLoadingFinish(s *network.Session, index int32, msg *protos.ReqPlayerLoadingFinish) {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	context.EventBus.Publish(events.PlayerLoadingFinish, player)
}

func (ps *PlayerRoute) ReqCreate(s *network.Session, msg *protos.ReqPlayerCreate) {
	if util.IsBlankString(msg.Name) {
		s.Send(&protos.ResPlayerCreate{Code: constants.I18N_COMMON_ILLEGAL_PARAMS}, 0)
		return
	}
	player := ps.service.Create(msg.Name, msg.Camp)
	s.Send(&protos.ResPlayerCreate{Code: 0, PlayerId: player.Id}, 0)
}

func (ps *PlayerRoute) ReqPlayerUpLevel(s *network.Session, index int32, msg *protos.ReqPlayerUpLevel) *protos.ResPlayerUpLevel {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoUpLevel(p, msg.ToLevel)
}

func (ps *PlayerRoute) ReqPlayerUpStage(s *network.Session, index int32, msg *protos.ReqPlayerUpStage) *protos.ResPlayerUpStage {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoUpStage(p)
}

func (ps *PlayerRoute) ReqPlayerRefreshScore(s *network.Session, index int32, msg *protos.ReqPlayerRefreshScore) *protos.ResPlayerRefreshScore {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	player.ClientScore = msg.Score
	context.EventBus.Publish(events.PlayerAttrChange, player)
	return &protos.ResPlayerRefreshScore{Code: 0}
}

func (ps *PlayerRoute) ReqEditClientData	(s *network.Session, index int32, msg *protos.ReqEditClientData) *protos.ResEditClientData {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	player.ClientData = msg.Data
	context.EventBus.Publish(events.PlayerAttrChange, player)
	return &protos.ResEditClientData{Code: 0}
}