package context

import (
	"github.com/forfun/gforgame/cache"
	"github.com/forfun/gforgame/common/eventbus"
	"github.com/forfun/gforgame/common/schedule"
	"github.com/forfun/gforgame/internal/infra/persistence"
	serverpkg "github.com/forfun/gforgame/network/server"
	"github.com/gin-gonic/gin"
)

var (
	CacheManager *cache.Manager
	DbService    *persistence.AsyncDBService
	GameServer   serverpkg.Server
	HttpServer   *gin.Engine
	EventBus      *eventbus.EventBus
	TaskScheduler schedule.TaskScheduler
)

func init() {
	CacheManager = cache.NewCacheManager()
	DbService = persistence.NewAsyncDbService()
	EventBus = eventbus.NewEventBus()
	TaskScheduler = schedule.NewDefaultTaskScheduler()
}
