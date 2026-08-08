package friend

import (
	"github.com/forfun/gforgame/cache"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"go.uber.org/dig"
)

// CachedFriendRepository 是好友仓储的缓存装饰器。
// 读：走 cache，miss 时由 cache loader 回源到 inner（MySQLFriendRepository）。
// 写：先持久化（inner），再更新缓存。
type CachedFriendRepository struct {
	inner *MySQLFriendRepository
	cache *cache.Manager
}

type CachedFriendRepositoryParams struct {
	dig.In
	Inner        *MySQLFriendRepository
	CacheManager *cache.Manager
}

func NewCachedFriendRepository(params CachedFriendRepositoryParams) *CachedFriendRepository {
	repo := &CachedFriendRepository{
		inner: params.Inner,
		cache: params.CacheManager,
	}
	// 构造时注册 cache 回源 loader，cache miss 时调 inner 从 db 加载。
	// 不依赖 InitServiceModules 反射调 Init：Services.FriendRepo 是接口类型字段，
	// 会被 InitServiceModules 的 reflect.Ptr 判断跳过，故 loader 必须在构造期注册。
	repo.cache.Register("friend", func(key string) (any, error) {
		return repo.inner.GetFriendEnt(key), nil
	})
	return repo
}

func (r *CachedFriendRepository) GetFriendEnt(playerId string) *playerdomain.Friend {
	c, err := r.cache.GetCache("friend")
	if err != nil {
		return nil
	}
	entity, err := c.Get(playerId)
	if err != nil || entity == nil {
		return nil
	}
	friend, _ := entity.(*playerdomain.Friend)
	return friend
}

func (r *CachedFriendRepository) SaveFriend(friend *playerdomain.Friend) {
	r.inner.SaveFriend(friend)
	if c, err := r.cache.GetCache("friend"); err == nil {
		c.Set(friend.Id, friend)
	}
}
