package chat

import (
	"time"

	"github.com/forfun/gforgame/internal/constants"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/idgen"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/protos"
	friendservice "github.com/forfun/gforgame/internal/service/friend"
	playerservice "github.com/forfun/gforgame/internal/service/player"
)

// 聊天模块
type ChatService struct {
	player   *playerservice.PlayerService
	handlers map[int32]ChatChannelHandler
}

func NewChatService(player *playerservice.PlayerService, friend *friendservice.FriendService) *ChatService {
	service := &ChatService{
		player:   player,
		handlers: make(map[int32]ChatChannelHandler),
	}
	service.handlers[constants.ChannelTypeFriend] = NewFriendChatChannelHandler(player, friend)
	service.handlers[constants.ChannelTypeWorld] = NewWorldChatChannelHandler(player)
	for _, chatHandler := range service.handlers {
		chatHandler.Init()
	}
	return service
}

func (s *ChatService) LoadOfflineMessages(player *playerdomain.Player) {
	offlineMessages := make([]*playerdomain.ChatMessage, 0)
	for _, handler := range s.handlers {
		offlineMessages = append(offlineMessages, handler.LoadOfflineMessages(player)...)
	}
	messages := make([]*protos.ChatMessageVo, 0)
	for _, message := range offlineMessages {
		sender := s.player.GetPlayerProfileById(message.SenderId)
		messageVo := &protos.ChatMessageVo{
			Id:         message.Id,
			Channel:    message.Channel,
			SenderId:   message.SenderId,
			ReceiverId: message.ReceiverId,
			Timestamp:  message.Timestamp,
			Content:    message.Content,
		}
		if sender != nil {
			messageVo.SenderName = sender.Name
			messageVo.SenderHead = sender.Head
		}
		messages = append(messages, messageVo)
	}
	push := &protos.PushChatNewMessage{
		Messages: messages,
	}
	io.NotifyPlayer(player, push)
}

func (s *ChatService) SendMessage(player *playerdomain.Player, msg *protos.ReqChat) *protos.ResChat {
	handler := s.handlers[int32(msg.Channel)]
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
	playerProfile := s.player.GetPlayerProfileById(player.Id)
	if playerProfile == nil {
		return &protos.ResChat{
			Code: -1,
		}
	}
	message := &playerdomain.ChatMessage{
		Id:         idgen.GetNextID(),
		Channel:    int32(msg.Channel),
		SenderId:   player.Id,
		ReceiverId: msg.Target,
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
