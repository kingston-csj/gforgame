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

// 常量定义
const (
	StatusStart  = 1 // 活动开始
	StatusClosed = 0 // 活动结束
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

type BaseActivityHandler struct {
	ActivityId    int32     // 活动ID
	LastEndDate   *time.Time // 最后结束日期
	StartTime     int64     // 活动开始时间戳（毫秒）
	EndTime       int64     // 活动结束时间戳（毫秒）
	Status       int                    // 活动状态
	CronMethods   map[string]reflect.Method // Cron 对应的方法映射
	ActivitySched *ActivityScheduler     // 依赖注入的调度器
	OnActivityStart func() error // 活动开始时的差异化逻辑
	OnActivityEnd   func() error // 活动结束时的差异化逻辑
}

// GetLastEndDate 获取最后结束日期
func (h *BaseActivityHandler) ActivityStart() {
	 h.Status = StatusStart
	 if h.OnActivityStart != nil {
		h.OnActivityStart()
	 }
}

// ActivityEnd 活动结束
func (h *BaseActivityHandler) ActivityEnd() {
	 h.Status = StatusClosed
	 if h.OnActivityEnd != nil {
		h.OnActivityEnd()
	 }
}

func (h *BaseActivityHandler) GetActivityId() int32 {
	 return h.ActivityId
}

type BaseHandlerGetter interface {
	GetBaseHandler() *BaseActivityHandler
}

// IsInActivity 接口方法的默认实现（子类可重写）
func (h *BaseActivityHandler) IsInActivity(player *player.Player) bool {
	// 通用逻辑：判断状态为开始，且当前时间在开始/结束时间范围内
	now := time.Now().UnixMilli()
	return h.Status == StatusStart && now >= h.StartTime && now <= h.EndTime
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
