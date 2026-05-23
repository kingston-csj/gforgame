package route

import (
	"github.com/forfun/gforgame/internal/context"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/protos"
	player "github.com/forfun/gforgame/internal/service/player"
	quest "github.com/forfun/gforgame/internal/service/quest"
	"github.com/forfun/gforgame/network"
)

type QuestRoute struct {
	network.Base
	service *quest.QuestService
	player  *player.PlayerService
}

func NewQuestRoute(service *quest.QuestService, playerService *player.PlayerService) *QuestRoute {
	return &QuestRoute{
		service: service,
		player:  playerService,
	}
}

func (c *QuestRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		c.service.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *QuestRoute) ReqTakeReward(s *network.Session, index int32, msg *protos.ReqQuestTakeReward) *protos.ResQuestTakeReward {
	player := ps.player.GetPlayerBySession(s)
	response, err := ps.service.TakeReward(player, msg.Id)
	if err != nil {
		return &protos.ResQuestTakeReward{
			Code: int32(err.Code()),
		}
	}
	return response	
}

func (ps *QuestRoute) ReqTakeProgressReward(s *network.Session, index int32, msg *protos.ReqQuestTakeProgressReward) *protos.ResQuestTakeProgressReward {
	player := ps.player.GetPlayerBySession(s)
	response := ps.service.TakeProgressReward(player, msg.Type)
	return response
}


func (ps *QuestRoute) ReqQuestEntrust(s *network.Session, index int32, msg *protos.ReqQuestEntrust) *protos.ResQuestEntrust {
	player := ps.player.GetPlayerBySession(s)
	err := ps.service.EntrustQuest(player, msg.QuestId, msg.HeroId)
	return &protos.ResQuestEntrust{
		Code: int32(err),
	}
}