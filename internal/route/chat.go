package route

import (
	"github.com/forfun/gforgame/internal/context"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/chat"
	playerservice "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/network"
)

type ChatRouter struct {
	network.Base
	service *chat.ChatService
}

func NewChatRoute() *ChatRouter {
	return &ChatRouter{
		service: chat.GetChatService(),
	}
}

func (rs *ChatRouter) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		rs.service.LoadOfflineMessages(data.(*playerdomain.Player))
	})
}

func (rs *ChatRouter) ReqChat(s *network.Session, index int32, msg *protos.ReqChat) *protos.ResChat {
	p := playerservice.GetPlayerService().GetPlayerBySession(s)
	response := rs.service.SendMessage(p, msg)
	return response
}
