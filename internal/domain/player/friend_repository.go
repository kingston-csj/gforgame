package player

// FriendRepository 是好友聚合根的仓储接口（依赖倒置锚点）。
// infra/repository/friend 提供实现，dig 装配时注入缓存装饰器。
type FriendRepository interface {
	// GetFriendEnt 按 playerId 加载好友聚合，不存在返回 nil。
	GetFriendEnt(playerId string) *Friend
	// SaveFriend 持久化好友聚合（含缓存更新）。
	SaveFriend(friend *Friend)
}
