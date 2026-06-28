package chat

import (
	util "github.com/forfun/gforgame/common/util/conv"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/protos"
	playerservice "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/network"
)

type ChatChannelHandler interface {

	// 初始化
	Init()

	// 检查是否可以发送
	CheckCanSend(player *playerdomain.Player, target string, content string) int

	// 保存到内存或者数据库
	SaveToDb(ChatMessage *playerdomain.ChatMessage)

	// 加载离线消息
	LoadOfflineMessages(player *playerdomain.Player) []*playerdomain.ChatMessage

	// 广播消息
	Broadcast(ChatMessage *playerdomain.ChatMessage)

	// 枚举消息接收者
	Receivers(message *playerdomain.ChatMessage) []string

	// 获取频道类型
	ChannelType() int32
}

type BaseChatChannelHandler struct {
	self   ChatChannelHandler // 指向子类
	player *playerservice.PlayerService
}

func NewBaseChatChannelHandler(self ChatChannelHandler, player *playerservice.PlayerService) *BaseChatChannelHandler {
	return &BaseChatChannelHandler{self: self, player: player}
}

func (b *BaseChatChannelHandler) Broadcast(message *playerdomain.ChatMessage) {
	// 过滤不在线的玩家
	receivers := b.self.Receivers(message)

	onlines := make([]*playerdomain.Player, 0)
	// 广播消息
	for _, receiver := range receivers {
		if network.IsOnline(receiver) {
			player := b.player.GetPlayer(receiver)
			if player != nil {
				onlines = append(onlines, player)
			}
		}
	}

	if len(onlines) <= 0 {
		return
	}

	sender := b.player.GetPlayerProfileById(message.SenderId)
	messageVo := &protos.ChatMessageVo{
		Id:         message.Id,
		Channel:    message.Channel,
		SenderId:   message.SenderId,
		SenderHead: sender.Head,
		SenderName: sender.Name,
		ReceiverId: message.ReceiverId,
		Timestamp:  message.Timestamp,
		Content:    message.Content,
	}

	if util.IsEmptyString(message.SenderId) {
		messageVo.SenderName = "系统"
	} else {
		playerProfile := b.player.GetPlayerProfileById(message.SenderId)
		if playerProfile != nil {
			messageVo.SenderName = playerProfile.Name
		}
	}

	push := &protos.PushChatNewMessage{
		Code:     0,
		Messages: make([]*protos.ChatMessageVo, 0),
	}
	push.Messages = append(push.Messages, messageVo)

	for _, p := range onlines {
		io.NotifyPlayer(p, push)
	}
}
