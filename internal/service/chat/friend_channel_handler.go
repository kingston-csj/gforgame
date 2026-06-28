package chat

import (
	"github.com/forfun/gforgame/internal/constants"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/service/friend"
	playerservice "github.com/forfun/gforgame/internal/service/player"
)

type FriendChannelHandler struct {
	*BaseChatChannelHandler
	friend *friend.FriendService
	player *playerservice.PlayerService
}

func (h *FriendChannelHandler) Init() {

}

func NewFriendChatChannelHandler(player *playerservice.PlayerService, friendService *friend.FriendService) *FriendChannelHandler {
	h := &FriendChannelHandler{
		friend: friendService,
		player: player,
	}
	h.BaseChatChannelHandler = NewBaseChatChannelHandler(h, player)
	return h
}

func (h *FriendChannelHandler) CheckCanSend(player *playerdomain.Player, target string, content string) int {
	if h.player.GetPlayerProfileById(target) == nil {
		return constants.I18N_COMMON_NOT_FOUND
	}
	if !h.friend.IsFriend(player.Id, target) {
		return constants.I18N_FRIEND_TIPS1
	}
	return 0
}

func (h *FriendChannelHandler) SaveToDb(message *playerdomain.ChatMessage) {
	from := h.player.GetPlayer(message.SenderId)
	from.ExtendBox.AddNewMessage(message)
	h.player.SavePlayer(from)

	to := h.player.GetPlayer(message.ReceiverId)
	to.ExtendBox.AddNewMessage(message)
	h.player.SavePlayer(to)
}

func (h *FriendChannelHandler) LoadOfflineMessages(player *playerdomain.Player) []*playerdomain.ChatMessage {
	offlineMessages := make([]*playerdomain.ChatMessage, 0)
	for _, v := range player.ExtendBox.PrivateChats {
		for _, v2 := range v {
			message := &playerdomain.ChatMessage{
				Id:         v2.Id,
				Channel:    v2.Channel,
				SenderId:   v2.SenderId,
				ReceiverId: v2.ReceiverId,
				SenderHead: v2.SenderHead,
				Timestamp:  v2.Timestamp,
				Content:    v2.Content,
			}
			offlineMessages = append(offlineMessages, message)
		}
	}
	return offlineMessages
}

func (h *FriendChannelHandler) Receivers(message *playerdomain.ChatMessage) []string {
	return []string{message.SenderId, message.ReceiverId}
}

func (h *FriendChannelHandler) ChannelType() int32 {
	return constants.ChannelTypeFriend
}
