package context

import (
	"io/github/gforgame/cache"
	"io/github/gforgame/eventbus"
	"io/github/gforgame/network/ws"

	mysqldb "io/github/gforgame/db"
)

var (
	CacheManager *cache.Manager
	DbService    *mysqldb.AsyncDbService
	// TcpServer    *tcp.TcpServer
	WsServer *ws.WsServer
	// HttpServer   *gin.Engine
	EventBus       *eventbus.EventBus
 
)

func init() {
	CacheManager = cache.NewCacheManager()
	DbService = mysqldb.NewAsyncDbService()
	EventBus = eventbus.NewEventBus()
}
