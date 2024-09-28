package context

import (
	"io/github/gforgame/cache"
)
import mysqldb "io/github/gforgame/db"

var (
	CacheManager *cache.Manager
	DbService    *mysqldb.AsyncDbService
)

func init() {
	CacheManager = cache.NewCacheManager()
	DbService = mysqldb.NewAsyncDbService()
}
