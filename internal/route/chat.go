package route

import (
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/chat"
	"github.com/forfun/gforgame/internal/service/player"
)


type ChatRouter struct {
	service *chat.ChatService
	player  *player.PlayerService
}

func NewChatRoute(service *chat.ChatService, playerService *player.PlayerService) *ChatRouter {
	return &ChatRouter{
		service: service,
		player:  playerService,
	}
}

func (rs *ChatRouter) ReqChat(playerId string, index int32, msg *protos.ReqChat) *protos.ResChat {
	p := rs.player.GetPlayer(playerId)
	response := rs.service.SendMessage(p, msg)
	return response
}
