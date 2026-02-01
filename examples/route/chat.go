package route

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/service/chat"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
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
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	response := rs.service.SendMessage(p, msg)
	return response
}
