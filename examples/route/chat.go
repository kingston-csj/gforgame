package route

import (
	"github.com/forfun/gforgame/examples/context"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	"github.com/forfun/gforgame/examples/events"
	"github.com/forfun/gforgame/examples/protos"
	"github.com/forfun/gforgame/examples/service/chat"
	playerservice "github.com/forfun/gforgame/examples/service/player"
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
