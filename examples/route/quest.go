package route

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	quest "io/github/gforgame/examples/service/quest"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
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
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	response, err := ps.service.TakeReward(player, msg.Id)
	if err != nil {
		return &protos.ResQuestTakeReward{
			Code: int32(err.Code()),
		}
	}
	return response	
}

func (ps *QuestRoute) ReqTakeProgressReward(s *network.Session, index int32, msg *protos.ReqQuestTakeProgressReward) *protos.ResQuestTakeProgressReward {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	response := ps.service.TakeProgressReward(player, msg.Type)
	return response
}


func (ps *QuestRoute) ReqQuestEntrust(s *network.Session, index int32, msg *protos.ReqQuestEntrust) *protos.ResQuestEntrust {
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	err := ps.service.EntrustQuest(player, msg.QuestId, msg.HeroId)
	return &protos.ResQuestEntrust{
		Code: int32(err),
	}
}