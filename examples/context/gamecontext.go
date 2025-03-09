package context

import (
	"io/github/gforgame/cache"
	"io/github/gforgame/network/ws"

	mysqldb "io/github/gforgame/db"
)

var (
	CacheManager *cache.Manager
	DbService    *mysqldb.AsyncDbService
	// TcpServer    *tcp.TcpServer
	WsServer     *ws.WsServer
	// HttpServer   *gin.Engine
)

func init() {
	CacheManager = cache.NewCacheManager()
	DbService = mysqldb.NewAsyncDbService()
}
