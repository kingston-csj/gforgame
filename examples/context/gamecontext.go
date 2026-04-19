package context

import (
	"io/github/gforgame/cache"
	"io/github/gforgame/common/eventbus"
	"io/github/gforgame/common/schedule"
	"io/github/gforgame/network/tcp"
	"io/github/gforgame/network/ws"

	"io/github/gforgame/examples/infra/persistence"
)

var (
	CacheManager *cache.Manager
	DbService    *persistence.AsyncDBService
	TcpServer    *tcp.TcpServer
	WsServer     *ws.WsServer
	// HttpServer   *gin.Engine
	EventBus      *eventbus.EventBus
	TaskScheduler schedule.TaskScheduler
)

func init() {
	CacheManager = cache.NewCacheManager()
	DbService = persistence.NewAsyncDbService()
	EventBus = eventbus.NewEventBus()
	TaskScheduler = schedule.NewDefaultTaskScheduler()
}
