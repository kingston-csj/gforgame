package route

import (
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/chat"
)

type ChatRouter struct {
	service    *chat.ChatService
	playerrepo *playerrepo.PlayerRepository
}

func NewChatRoute(service *chat.ChatService, playerRepo *playerrepo.PlayerRepository) *ChatRouter {
	return &ChatRouter{
		service:    service,
		playerrepo: playerRepo,
	}
}

func (rs *ChatRouter) ReqChat(playerId string, index int32, msg *protos.ReqChat) *protos.ResChat {
	p := rs.playerrepo.GetPlayer(playerId)
	response := rs.service.SendMessage(p, msg)
	return response
}
