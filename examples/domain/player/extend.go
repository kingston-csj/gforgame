package player

type ExtendBox struct {

	// 私聊消息 key为对方id，value为消息列表
	PrivateChats map[string][]ChatMessage
	// vip每个周期的充值金额
	VipPeriodMoney float32 `json:"vipPeriodMoney"`
	// vip过期时间
	VipExpiredTime int64 `json:"vipExpiredTime"`
	// 成就积分
	AchievementScore int32 `json:"achievementScore"`
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

func (b *ExtendBox) AddVipPeriodMoney(money float32) float32 {
	b.VipPeriodMoney += money
	return b.VipPeriodMoney
}
