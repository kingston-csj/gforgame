package chat

import (
	"io/github/gforgame/examples/constants"
	playerdomain "io/github/gforgame/examples/domain/player"
	friendservice "io/github/gforgame/examples/friend"
	playerservice "io/github/gforgame/examples/service/player"
)

type FriendChannelHandler struct {
	BaseChatChannelHandler
}

func (h *FriendChannelHandler) Init() {

}

func (h *FriendChannelHandler) CheckCanSend(player *playerdomain.Player, target string, content string) int {
	if playerservice.GetPlayerService().GetPlayerProfileById(target) == nil {
		return constants.I18N_COMMON_NOT_FOUND
	}
	if friendservice.GetFriendService().IsFriend(player.Id, target) {
		return constants.I18N_FRIEND_TIPS1
	}
	return 0
}

func (h *FriendChannelHandler) SaveToDb(message *playerdomain.ChatMessage) {
	from := playerservice.GetPlayerService().GetPlayer(message.SenderId)
	from.ExtendBox.AddNewMessage(message)
	playerservice.GetPlayerService().SavePlayer(from)

	to := playerservice.GetPlayerService().GetPlayer(message.ReceiverId)
	to.ExtendBox.AddNewMessage(message)
	playerservice.GetPlayerService().SavePlayer(to)
}

func (h *FriendChannelHandler) LoadOfflineMessages(player *playerdomain.Player) []*playerdomain.ChatMessage {
	offlineMessages := make([]*playerdomain.ChatMessage, 0)
	for _, v := range player.ExtendBox.PrivateChats {
		for _, v2 := range v {
			message := &playerdomain.ChatMessage{
				Id:         v2.Id,
				Channel:    v2.Channel,
				SenderId:   v2.SenderId,
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
