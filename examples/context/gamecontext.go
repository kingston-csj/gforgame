package context

import "io/github/gforgame/cache"

var (
	CacheManager *cache.CacheManager
)

func init() {

	CacheManager = cache.NewCacheManager()
}
