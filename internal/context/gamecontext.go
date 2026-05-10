package context

import (
	"github.com/forfun/gforgame/cache"
	"github.com/forfun/gforgame/common/eventbus"
	"github.com/forfun/gforgame/common/schedule"
	"github.com/forfun/gforgame/network/tcp"
	"github.com/forfun/gforgame/network/ws"

	"github.com/forfun/gforgame/internal/infra/persistence"
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
