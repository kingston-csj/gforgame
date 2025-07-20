package player

type ExtendBox struct {

	// 私聊消息 key为对方id，value为消息列表
	PrivateChats map[string][]ChatMessage
}

func (b *ExtendBox) AddNewMessage(message *ChatMessage) {
	if b.PrivateChats == nil {
		b.PrivateChats = make(map[string][]ChatMessage)
	}
	chatMessages := b.PrivateChats[message.ReceiverId]
	if chatMessages == nil {
		chatMessages = make([]ChatMessage, 0)
	}
	chatMessages = append(chatMessages, *message)
	b.PrivateChats[message.ReceiverId] = chatMessages
}
