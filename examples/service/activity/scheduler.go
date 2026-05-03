package activity

import (
	"fmt"
	"io/github/gforgame/common/logger"
	"io/github/gforgame/common/schedule"
	"io/github/gforgame/common/util/timeutil"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	"reflect"
	"sync"
	"time"
)

// innerTaskMap 内层任务映射，保证并发安
type innerTaskMap struct {
	mu    sync.RWMutex
	tasks map[string]schedule.Cancellable
}


// ActivityScheduler 活动调度
type ActivityScheduler struct {
	// taskScheduler 定时任务调度
	taskScheduler schedule.TaskScheduler
	// scheduledTasks 记录活动调度任务，双层并发安全映
	// key: 活动ID（int32 value: *innerTaskMap（状态名-> 可取消任务）
	scheduledTasks sync.Map
}

// NewActivityScheduler 创建新的活动调度
func NewActivityScheduler(taskScheduler schedule.TaskScheduler) *ActivityScheduler {
	return &ActivityScheduler{
		taskScheduler: taskScheduler,
		scheduledTasks: sync.Map{},
	}
}

// ScheduleActivity 调度活动
func (a *ActivityScheduler) ScheduleActivity(handler ActivityHandler) {
	if handler == nil {
		return
	}
	getter, ok := handler.(BaseHandlerGetter)
	if !ok {
		activityId := handler.GetActivityId()
		logger.Info("ScheduleActivity: handler does not implement BaseHandlerGetter, activityId= " + string(activityId))
		return
	}
	activityId := getter.GetBaseHandler().GetActivityId()

	defer func() {
		// 捕获 panic
		if r := recover(); r != nil {
			logger.Error("ScheduleActivity: panic recovered: %v", fmt.Errorf("%v", r))
		}
	}()

	// 取消已存在的调度
	a.CancelScheduledActivity(activityId)

	// 获取活动数据和时间信
	activityData  := config.QueryById[configdomain.ActivityData](activityId)
	startCron := activityData.Start
	endCron := activityData.End

	// 获取参考日期和当前时间
	referDate := a.getReferenceDate(getter, startCron)
	now := time.Now()

	// 计算开始和结束时间
	startDate, _ := schedule.GetNextTriggerTimeAfter(startCron, referDate)
	endDate, _ := schedule.GetNextTriggerTimeAfter(endCron, referDate)

	// 修正：启服时间超过开始时间，根据周期时长计算开始时
	if startDate.After(endDate) {
		realEndTime := endDate.UnixMilli()
		nextStartTime := startDate.UnixMilli()
		nextEndTime, _ := schedule.GetNextTriggerTimeAfter(endCron, startDate)
		nextEndTimeMilli := nextEndTime.UnixMilli()

		// 活动周期
		duration := nextEndTimeMilli - nextStartTime
		realStartTime := realEndTime - duration
		startDate = time.UnixMilli(realStartTime)
	}

	// 设置活动结束时间
	getter.GetBaseHandler().EndTime = endDate.UnixMilli()

	// 如果结束时间在当前时间之后，调度相关任务
	if endDate.After(now) {
		a.scheduleStartTask(getter, startDate, now, endDate)
		a.scheduleEndTask(activityId, endDate, now)
		a.scheduleCronTasks(getter, activityId, now)
	}
}

// getReferenceDate 获取参考日
func (a *ActivityScheduler) getReferenceDate(handler BaseHandlerGetter, startCron string) time.Time {
	// 程序运行期间，参考日期以最近的活动结束时间为准
	if handler.GetBaseHandler().LastEndDate != nil {
		return *handler.GetBaseHandler().LastEndDate
	}

	// 周期性活动，取当前时
	if schedule.GetParser(startCron).IsPeriodicExpression(startCron) {
		return time.Now()
	}

	// 刚启服场景，1970-01-01 作为参考日
	defaultDate, _ := timeutil.ParseLocalTime("1970-01-01 00:00:00")
	return defaultDate
}

// scheduleStartTask 调度活动开始任
func (a *ActivityScheduler) scheduleStartTask(
	handler BaseHandlerGetter,
	startDate, now, endDate time.Time) {
	// 设置活动结束时间
	handler.GetBaseHandler().EndTime = endDate.UnixMilli()
	activityId := handler.GetBaseHandler().ActivityId	
	// 定义开始任
	startTask := func() {
		logger.Log(constants.LoggerActivity, "type", "start", "activityId", activityId)
		// 调用活动开始方
		handler.GetBaseHandler().ActivityStart()
	}

	// 开始时间已过期，立即执
	if startDate.Before(now) {
		go startTask() // 协程执行，避免阻
		return
	}

	// 调度未来执行的任
	future, err := a.taskScheduler.Schedule(startTask, startDate.Sub(now).Milliseconds())
	if err != nil {
		logger.Error("", err)
		return
	}

	// 注册任务
	a.registerTask(activityId, "start", future)
}

// scheduleEndTask 调度活动结束任务（对Java scheduleEndTask 方法
func (a *ActivityScheduler) scheduleEndTask(activityId int32, endDate, now time.Time) {
	// 定义结束任务
	endTask := func() {
		logger.Log(constants.LoggerActivity, "type", "end", "activityId", activityId)
		// 调用活动结束方法
		handler, err := GetHandler(activityId)
		if err != nil {
			logger.Error("", err)
			return
		}
		handler.ActivityEnd()
	}

	// 结束时间在当前时间之后，调度任务
	if endDate.After(now) {
		future, err := a.taskScheduler.Schedule(endTask, endDate.Sub(now).Milliseconds())
		if err != nil {
			logger.Error("", err)
			return
		}
		// 注册任务
		a.registerTask(activityId, "end", future)
	}
}

// scheduleCronTasks 调度其他特殊时间任务（对Java scheduleCronTasks 方法
func (a *ActivityScheduler) scheduleCronTasks(handler BaseHandlerGetter, activityId int32, now time.Time) {
	// 遍历 Cron 方法映射
	for cronExpr, method := range handler.GetBaseHandler().CronMethods {
		cronDate, _ := schedule.GetNextTriggerTimeAfter(cronExpr, now)
		if !cronDate.After(now) {
			continue
		}

		// 定义 Cron 任务
		cronTask := func() {
			logger.Log(constants.LoggerActivity, "type", method.Name, "activityId", activityId)
			// 调用对应方法
			handler, err := GetHandler(activityId)
			if err != nil {
				logger.Error("", err)
				return
			}

			// 反射调用方法（对Java method.invoke
			handlerValue := reflect.ValueOf(handler)
			methodValue := handlerValue.MethodByName(method.Name)
			if !methodValue.IsValid() {
				logger.Info("activity handler method not found")
				return
			}
			// 调用无参方法
			methodValue.Call(nil)
		}

		// 调度 Cron 任务
		future, err := a.taskScheduler.Schedule(cronTask, cronDate.Sub(now).Milliseconds())
		if err != nil {
			logger.Error("", err)
			continue
		}

		// 注册任务
		a.registerTask(activityId, cronExpr, future)
	}
}

// registerTask 注册任务（对Java registerTask 方法
func (a *ActivityScheduler) registerTask(activityId int32, name string, task schedule.Cancellable) {
	// 外层 sync.Map 获取内层任务映射
	innerVal, ok := a.scheduledTasks.Load(activityId)
	var innerMap *innerTaskMap
	if !ok {
		// 不存在则创建新的内层任务映射
		innerMap = &innerTaskMap{
			tasks: make(map[string]schedule.Cancellable),
		}
		a.scheduledTasks.Store(activityId, innerMap)
	} else {
		innerMap = innerVal.(*innerTaskMap)
	}

	// 内层加锁更新任务
	innerMap.mu.Lock()
	defer innerMap.mu.Unlock()

	// 如有旧任务，先取
	if oldTask, exists := innerMap.tasks[name]; exists {
		oldTask.Cancel()
	}
	innerMap.tasks[name] = task
}

// CancelScheduledActivity 取消活动调度（对Java cancelScheduledActivity 方法
func (a *ActivityScheduler) CancelScheduledActivity(activityId int32) {
	// 从外sync.Map 移除内层任务映射
	innerVal, ok := a.scheduledTasks.LoadAndDelete(activityId)
	if !ok {
		return
	}

	innerMap := innerVal.(*innerTaskMap)
	innerMap.mu.RLock()
	defer innerMap.mu.RUnlock()

	// 取消所有任
	for _, task := range innerMap.tasks {
		task.Cancel()
	}
}
