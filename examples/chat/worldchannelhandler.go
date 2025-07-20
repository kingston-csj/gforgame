package chat

import (
	"io/github/gforgame/ds/list"
	"io/github/gforgame/examples/constants"
	playerdomain "io/github/gforgame/examples/domain/player"
	network "io/github/gforgame/network"
	"time"
)

type WorldChannelHandler struct {
	BaseChatChannelHandler
	MsgQueue *list.LimitedList[*playerdomain.ChatMessage]
}

func (h *WorldChannelHandler) Init() {
	h.MsgQueue = list.NewLimitedList[*playerdomain.ChatMessage](100)

	msg := &playerdomain.ChatMessage{
		Id:         "1",
		Channel:    constants.ChannelTypeWorld,
		SenderId:   "",
		SenderHead: 0,
		Timestamp:  time.Now().Unix(),
		Content:    "亲爱的玩家，欢迎回来！请友好交流，共建美好聊天环境！",
	}

	h.MsgQueue.Push(msg)
}

func (h *WorldChannelHandler) CheckCanSend(player *playerdomain.Player, target string, content string) int {
	return 0
}

func (h *WorldChannelHandler) SaveToDb(message *playerdomain.ChatMessage) {
	h.MsgQueue.Push(message)
}

func (h *WorldChannelHandler) LoadOfflineMessages(player *playerdomain.Player) []*playerdomain.ChatMessage {
	offlineMessages := make([]*playerdomain.ChatMessage, 0)
	h.MsgQueue.Each(func(msg *playerdomain.ChatMessage) {
		offlineMessages = append(offlineMessages, msg)
	})
	return offlineMessages
}

func (h *WorldChannelHandler) Receivers(message *playerdomain.ChatMessage) []string {
	onlinePlayerIds := network.GetAllOnlinePlayerIds()
	receivers := make([]string, 0)
	for _, playerId := range onlinePlayerIds {
		receivers = append(receivers, playerId)
	}
	return receivers
}

func (h *WorldChannelHandler) ChannelType() int32 {
	return constants.ChannelTypeWorld
}
