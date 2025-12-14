package chat

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type ChatController struct {
	network.Base
}

func NewChatController() *ChatController {
	return &ChatController{}
}

func (rs *ChatController) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		GetChatService().LoadOfflineMessages(data.(*playerdomain.Player))
	})
}

func (rs *ChatController) ReqChat(s *network.Session, index int, msg *protos.ReqChat) *protos.ResChat {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	response := GetChatService().SendMessage(p, msg)
	return response
}
