package route

import (
	"github.com/forfun/gforgame/common/util"
	"github.com/forfun/gforgame/internal/constants"
	"github.com/forfun/gforgame/internal/context"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	mysqldb "github.com/forfun/gforgame/internal/infra/persistence"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/player"
	playerservice "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/network"

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
		p.AfterLoad()
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
				player := playerservice.GetPlayerService().GetPlayerBySession(s)
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
	player := playerservice.GetPlayerService().GetPlayerBySession(s)
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
	p := playerservice.GetPlayerService().GetPlayerBySession(s)
	return ps.service.DoUpLevel(p, msg.ToLevel)
}

func (ps *PlayerRoute) ReqPlayerUpStage(s *network.Session, index int32, msg *protos.ReqPlayerUpStage) *protos.ResPlayerUpStage {
	p := playerservice.GetPlayerService().GetPlayerBySession(s)
	return ps.service.DoUpStage(p)
}

func (ps *PlayerRoute) ReqPlayerRefreshScore(s *network.Session, index int32, msg *protos.ReqPlayerRefreshScore) *protos.ResPlayerRefreshScore {
	player := playerservice.GetPlayerService().GetPlayerBySession(s)
	player.ClientScore = msg.Score
	context.EventBus.Publish(events.PlayerAttrChange, player)
	return &protos.ResPlayerRefreshScore{Code: 0}
}

func (ps *PlayerRoute) ReqEditClientData	(s *network.Session, index int32, msg *protos.ReqEditClientData) *protos.ResEditClientData {
	player := playerservice.GetPlayerService().GetPlayerBySession(s)
	player.ClientData = msg.Data
	context.EventBus.Publish(events.PlayerAttrChange, player)
	return &protos.ResEditClientData{Code: 0}
}