package chat

import (
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
	"sync"

	playerdomain "io/github/gforgame/examples/domain/player"
	playerservice "io/github/gforgame/examples/service/player"
	"io/github/gforgame/util"
	"time"
)

type ChatService struct {
}

var (
	instance            *ChatService
	once                sync.Once
	ChatChannelHandlers map[int32]ChatChannelHandler = make(map[int32]ChatChannelHandler)
)

func GetChatService() *ChatService {
	once.Do(func() {
		instance = &ChatService{}
		ChatChannelHandlers[constants.ChannelTypeFriend] = &FriendChannelHandler{}
		ChatChannelHandlers[constants.ChannelTypeWorld] = &WorldChannelHandler{}
		for _, handler := range ChatChannelHandlers {
			handler.Init()
		}
	})
	return instance

}

func (s *ChatService) LoadOfflineMessages(player *playerdomain.Player) {
	offlineMessages := make([]*playerdomain.ChatMessage, 0)
	for _, handler := range ChatChannelHandlers {
		offlineMessages = append(offlineMessages, handler.LoadOfflineMessages(player)...)
	}
	messages := make([]*protos.ChatMessageVo, 0)
	for _, message := range offlineMessages {
		messages = append(messages, &protos.ChatMessageVo{
			Id:         message.Id,
			Channel:    message.Channel,
			SenderId:   message.SenderId,
			SenderHead: message.SenderHead,
			Timestamp:  message.Timestamp,
			Content:    message.Content,
		})
	}
	push := &protos.PushChatNewMessage{
		Messages: messages,
	}
	io.NotifyPlayer(player, push)
}

func (s *ChatService) SendMessage(player *playerdomain.Player, msg *protos.ReqChat) *protos.ResChat {
	handler := ChatChannelHandlers[int32(msg.Channel)]
	if handler == nil {
		return &protos.ResChat{
			Code: -1,
		}
	}
	code := handler.CheckCanSend(player, msg.Target, msg.Content)
	if code != 0 {
		return &protos.ResChat{
			Code: code,
		}
	}
	playerProfile := playerservice.GetPlayerService().GetPlayerProfileById(player.Id)
	if playerProfile == nil {
		return &protos.ResChat{
			Code: -1,
		}
	}
	message := &playerdomain.ChatMessage{
		Id:         util.GetNextID(),
		Channel:    int32(msg.Channel),
		SenderId:   player.Id,
		SenderHead: playerProfile.Head,
		Timestamp:  time.Now().Unix(),
		Content:    msg.Content,
	}
	handler.SaveToDb(message)
	handler.Broadcast(message)
	return &protos.ResChat{
		Code: code,
	}
}
