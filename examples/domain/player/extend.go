package player

type ExtendBox struct {

	// 私聊消息 key为对方id，value为消息列表
	PrivateChats map[string][]ChatMessage
	// vip每个周期的充值金额
	VipPeriodMoney float32
	// vip过期时间
	VipExpiredTime int64
	// 成就积分
	AchievementScore int32
	// 客户端事件统计
	ClientEvents map[int32]int32
	// 材料图鉴
	ItemCatalogModel CatalogModel
	// 店铺图鉴
	SitemCatalogModel CatalogModel
	// 菜单图鉴
	MenuCatalogModel CatalogModel
}

func (b *ExtendBox) AfterLoad() {
	if b.PrivateChats == nil {
		b.PrivateChats = make(map[string][]ChatMessage)
	}
	if b.ClientEvents == nil {
		b.ClientEvents = make(map[int32]int32)
	}
	b.ItemCatalogModel.AfterLoad()
	b.SitemCatalogModel.AfterLoad()
	b.MenuCatalogModel.AfterLoad()
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
