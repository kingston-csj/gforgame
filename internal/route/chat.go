package route

import (
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/chat"
	"github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/network"
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

func (rs *ChatRouter) ReqChat(s *network.Session, index int32, msg *protos.ReqChat) *protos.ResChat {
	p := rs.player.GetPlayerBySession(s)
	response := rs.service.SendMessage(p, msg)
	return response
}
