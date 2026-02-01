package chat

import (
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	playerservice "io/github/gforgame/examples/service/player"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"
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

	// 获取消息
	Receivers(message *playerdomain.ChatMessage) []string

	// 获取频道类型
	ChannelType() int32
}

type BaseChatChannelHandler struct {
	ChatChannelHandler
}

func (b *BaseChatChannelHandler) Broadcast(message *playerdomain.ChatMessage) {
	// 过滤不在线的玩家
	receivers := b.Receivers(message)

	onlines := make([]*playerdomain.Player, 0)
	// 广播消息
	for _, receiver := range receivers {
		if network.IsOnline(receiver) {
			player := network.GetPlayerByPlayerId(receiver)
			if player != nil {
				onlines = append(onlines, player.(*playerdomain.Player))
			}
		}
	}

	if len(onlines) <= 0 {
		return
	}

	messageVo := &protos.ChatMessageVo{
		Id:         message.Id,
		Channel:    message.Channel,
		SenderId:   message.SenderId,
		SenderHead: message.SenderHead,
		Timestamp:  message.Timestamp,
		Content:    message.Content,
	}

	if util.IsEmptyString(message.SenderId) {
		messageVo.SenderName = "系统"
	} else {
		playerProfile := playerservice.GetPlayerService().GetPlayerProfileById(message.SenderId)
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
