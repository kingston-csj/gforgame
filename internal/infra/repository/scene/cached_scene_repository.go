package scene

import (
	"github.com/forfun/gforgame/cache"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"go.uber.org/dig"
)

// CachedSceneRepository 是场景仓储的缓存装饰器。
// 读：走 cache，miss 时由 cache loader 回源到 inner（MySQLSceneRepository）。
// 写：先持久化（inner），再更新缓存。
type CachedSceneRepository struct {
	inner *MySQLSceneRepository
	cache *cache.Manager
}

type CachedSceneRepositoryParams struct {
	dig.In
	Inner        *MySQLSceneRepository
	CacheManager *cache.Manager
}

func NewCachedSceneRepository(params CachedSceneRepositoryParams) *CachedSceneRepository {
	repo := &CachedSceneRepository{
		inner: params.Inner,
		cache: params.CacheManager,
	}
	// 构造时注册 cache 回源 loader，cache miss 时调 inner 从 db 加载。
	repo.cache.Register("scene", func(key string) (any, error) {
		return repo.inner.GetSceneByID(key), nil
	})
	return repo
}

func (r *CachedSceneRepository) GetScene(key string) *playerdomain.Scene {
	c, err := r.cache.GetCache("scene")
	if err != nil {
		return nil
	}
	entity, err := c.Get(key)
	if err != nil || entity == nil {
		return nil
	}
	scene, _ := entity.(*playerdomain.Scene)
	return scene
}

func (r *CachedSceneRepository) SaveScene(scene *playerdomain.Scene) {
	r.inner.SaveSceneToDb(scene)
	if c, err := r.cache.GetCache("scene"); err == nil {
		c.Set(scene.GetId(), scene)
	}
}
