package route

import (
	"github.com/forfun/gforgame/examples/context"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	"github.com/forfun/gforgame/examples/events"
	"github.com/forfun/gforgame/examples/protos"
	playerservice "github.com/forfun/gforgame/examples/service/player"
	quest "github.com/forfun/gforgame/examples/service/quest"
	"github.com/forfun/gforgame/network"
)

type QuestRoute struct {
	network.Base
	service *quest.QuestService
}

func NewQuestRoute() *QuestRoute {
	return &QuestRoute{
		service: quest.GetQuestService(),
	}
}

func (c *QuestRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		c.service.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *QuestRoute) ReqTakeReward(s *network.Session, index int32, msg *protos.ReqQuestTakeReward) *protos.ResQuestTakeReward {
	player := playerservice.GetPlayerService().GetPlayerBySession(s)
	response, err := ps.service.TakeReward(player, msg.Id)
	if err != nil {
		return &protos.ResQuestTakeReward{
			Code: int32(err.Code()),
		}
	}
	return response	
}

func (ps *QuestRoute) ReqTakeProgressReward(s *network.Session, index int32, msg *protos.ReqQuestTakeProgressReward) *protos.ResQuestTakeProgressReward {
	player := playerservice.GetPlayerService().GetPlayerBySession(s)
	response := ps.service.TakeProgressReward(player, msg.Type)
	return response
}


func (ps *QuestRoute) ReqQuestEntrust(s *network.Session, index int32, msg *protos.ReqQuestEntrust) *protos.ResQuestEntrust {
	player := playerservice.GetPlayerService().GetPlayerBySession(s)
	err := ps.service.EntrustQuest(player, msg.QuestId, msg.HeroId)
	return &protos.ResQuestEntrust{
		Code: int32(err),
	}
}