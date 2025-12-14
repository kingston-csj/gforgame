package activity

import (
	"errors"
	"io/github/gforgame/examples/domain/player"
	"reflect"
	"time"
)

var (
	// handlers 存储活动ID到处理程序的映射
	handlers map[int32]ActivityHandler = make(map[int32]ActivityHandler)
)

type ActivityHandler interface {

	// LoadActivityInfo 加载玩家活动信息
	LoadActivityInfo(player *player.Player)

	// IsInActivity 判断玩家是否在活动中
	// 同时判断是否在活动期间，以及是否有参与资格
	IsInActivity(player *player.Player) bool

	// GetActivityId 获取当前活动ID
	GetActivityId() int32

	// ActivityStart 活动开始
	ActivityStart()

	// ActivityEnd 活动结束
	ActivityEnd()
}

type AbsActivityHandler struct {
	ActivityId    int32     // 活动ID
	LastEndDate   *time.Time // 最后结束日期
	EndTime       int64     // 活动结束时间戳（毫秒）
	CronMethods   map[string]reflect.Method // Cron 对应的方法映射
}


// ActivityData 活动数据结构体
type ActivityData struct {
	ActivityId int32  // 活动ID
	Start      string // 开始时间 Cron 表达式
	End        string // 结束时间 Cron 表达式
}

func registerHandler(activityId int32, handler ActivityHandler) {
	handlers[activityId] = handler
}

func GetHandler(activityId int32) (ActivityHandler, error) {
	handler, ok := handlers[activityId]
	if !ok {
		return nil, errors.New("activity handler not found")
	}
	return handler, nil
}