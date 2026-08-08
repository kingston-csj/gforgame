package route

import (
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/protos"
	quest "github.com/forfun/gforgame/internal/service/quest"
)

type QuestRoute struct {
	service *quest.QuestService
	repo    *playerrepo.PlayerRepository
}

func NewQuestRoute(service *quest.QuestService, repo *playerrepo.PlayerRepository) *QuestRoute {
	return &QuestRoute{
		service: service,
		repo:    repo,
	}
}

func (ps *QuestRoute) ReqTakeReward(playerId string, index int32, msg *protos.ReqQuestTakeReward) *protos.ResQuestTakeReward {
	player := ps.repo.GetPlayer(playerId)
	response, err := ps.service.TakeReward(player, msg.Id)
	if err != nil {
		return &protos.ResQuestTakeReward{
			Code: int32(err.Code()),
		}
	}
	return response
}

func (ps *QuestRoute) ReqTakeProgressReward(playerId string, index int32, msg *protos.ReqQuestTakeProgressReward) *protos.ResQuestTakeProgressReward {
	player := ps.repo.GetPlayer(playerId)
	response := ps.service.TakeProgressReward(player, msg.Type)
	return response
}

func (ps *QuestRoute) ReqQuestEntrust(playerId string, index int32, msg *protos.ReqQuestEntrust) *protos.ResQuestEntrust {
	player := ps.repo.GetPlayer(playerId)
	err := ps.service.EntrustQuest(player, msg.QuestId, msg.HeroId)
	return &protos.ResQuestEntrust{
		Code: int32(err),
	}
}
