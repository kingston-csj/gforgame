package route

import (
	"github.com/forfun/gforgame/internal/context"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/chat"
	"github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/network"
)


type ChatRouter struct {
	network.Base
	service *chat.ChatService
	player  *player.PlayerService
}

func NewChatRoute(service *chat.ChatService, playerService *player.PlayerService) *ChatRouter {
	return &ChatRouter{
		service: service,
		player:  playerService,
	}
}

func (rs *ChatRouter) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		rs.service.LoadOfflineMessages(data.(*playerdomain.Player))
	})
}

func (rs *ChatRouter) ReqChat(s *network.Session, index int32, msg *protos.ReqChat) *protos.ResChat {
	p := rs.player.GetPlayerBySession(s)
	response := rs.service.SendMessage(p, msg)
	return response
}
