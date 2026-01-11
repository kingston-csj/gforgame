package reward

import "io/github/gforgame/examples/domain/player"

type Reward interface {

	// 返回奖励数量
	GetAmount() int

	// 增加奖励数量
	AddAmount(amount int)
	
	// 验证能发发放，例如背包已满不可发放
	// 若验证不通过则返回错误
	Verify(player *player.Player) error

	// 验证能发发放，例如背包已满不可发放
	VerifySliently(player *player.Player) bool

	// 真正的发放逻辑
	Reward(player *player.Player, actionType int)

	GetType() string

	// 序列化为字符串，供客户端解析
	Serial() string

}
