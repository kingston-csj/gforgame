package context

import (
	"io/github/gforgame/cache"
	"io/github/gforgame/network/tcp"

	mysqldb "io/github/gforgame/db"

	"github.com/gin-gonic/gin"
)

var (
	CacheManager *cache.Manager
	DbService    *mysqldb.AsyncDbService
	TcpServer    *tcp.TcpServer
	HttpServer   *gin.Engine
)

func init() {
	CacheManager = cache.NewCacheManager()
	DbService = mysqldb.NewAsyncDbService()
}
