package activity

import (
	"io/github/gforgame/examples/config"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/logger"
	"io/github/gforgame/schedule"
	"io/github/gforgame/util/timeutil"
	"reflect"
	"sync"
	"time"
)


type Cancellable interface {
	// Cancel 取消任务
	Cancel() bool
}

// TaskScheduler 定时任务调度接口
type TaskScheduler interface {
	// Schedule 调度一次性任务，在指定时间执行
	Schedule(task func(), execTime time.Time) (Cancellable, error)
}


// innerTaskMap 内层任务映射，保证并发安全
type innerTaskMap struct {
	mu    sync.RWMutex
	tasks map[string]Cancellable
}


// ActivityScheduler 活动调度器
type ActivityScheduler struct {
	// taskScheduler 定时任务调度器 
	taskScheduler TaskScheduler

	// scheduledTasks 记录活动调度任务，双层并发安全映射
	// key: 活动ID（int32）, value: *innerTaskMap（状态名称 -> 可取消任务）
	scheduledTasks sync.Map
}

// NewActivityScheduler 创建新的活动调度器
func NewActivityScheduler(taskScheduler TaskScheduler) *ActivityScheduler {
	return &ActivityScheduler{
		taskScheduler: taskScheduler,
		scheduledTasks: sync.Map{},
	}
}

// ScheduleActivity 调度活动
func (a *ActivityScheduler) ScheduleActivity(activity *AbsActivityHandler) {
	if activity == nil {
		return
	}
	activityId := activity.ActivityId

	defer func() {
		// 捕获 panic（对应 Java 的 catch 异常）
		if r := recover(); r != nil {
			logger.Error(  r.(error))
		}
	}()

	// 取消已存在的调度
	a.CancelScheduledActivity(activityId)

	// 获取活动数据和时间信息
	activityData  := config.QueryById[configdomain.ActivityData](activityId)
	startCron := activityData.Start
	endCron := activityData.End

	// 获取参考日期和当前时间
	referDate := a.getReferenceDate(activity, startCron)
	now := time.Now()

	// 计算开始和结束时间
	startDate, _ := schedule.GetNextTriggerTimeAfter(startCron, referDate)
	endDate, _ := schedule.GetNextTriggerTimeAfter(endCron, referDate)

	// 修正：启服时间超过开始时间，根据周期时长计算开始时间
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

	// 设置活动结束时间戳
	activity.EndTime = endDate.UnixMilli()

	// 如果结束时间在当前时间之后，调度相关任务
	if endDate.After(now) {
		a.scheduleStartTask(activity, activityId, startDate, now, endDate)
		a.scheduleEndTask(activityId, endDate, now)
		a.scheduleCronTasks(activity, activityId, now)
	}
}

// getReferenceDate 获取参考日期（对应 Java 的 getReferenceDate 方法）
func (a *ActivityScheduler) getReferenceDate(activity *AbsActivityHandler, startCron string) time.Time {
	// 程序运行期间，参考日期以最近的活动结束时间为准
	if activity.LastEndDate != nil {
		return *activity.LastEndDate
	}

	// 周期性活动，取当前时间
	if schedule.GetParser(startCron).IsPeriodicExpression(startCron) {
		return time.Now()
	}

	// 刚启服场景，取 1970-01-01 作为参考日期
	defaultDate, _ := timeutil.ParseLocalTime("1970-01-01 00:00:00")
	return defaultDate
}

// scheduleStartTask 调度活动开始任务（对应 Java 的 scheduleStartTask 方法）
func (a *ActivityScheduler) scheduleStartTask(
	activity *AbsActivityHandler,
	activityId int32,
	startDate, now, endDate time.Time) {
	// 设置活动结束时间
	activity.EndTime = endDate.UnixMilli()

	// 定义开始任务
	startTask := func() {
		logger.Log(logger.Activity, "type", "start", "activityId", activityId)
		// 调用活动开始方法
		handler, err := GetHandler(activityId)
		if err != nil {
			logger.Error(err)
			return
		}
		handler.ActivityStart()
	}

	// 开始时间已过期，立即执行
	if startDate.Before(now) {
		go startTask() // 协程执行，避免阻塞
		return
	}

	// 调度未来执行的任务
	future, err := a.taskScheduler.Schedule(startTask, startDate)
	if err != nil {
		logger.Error(err)
		return
	}

	// 注册任务
	a.registerTask(activityId, "start", future)
}

// scheduleEndTask 调度活动结束任务（对应 Java 的 scheduleEndTask 方法）
func (a *ActivityScheduler) scheduleEndTask(activityId int32, endDate, now time.Time) {
	// 定义结束任务
	endTask := func() {
		logger.Log(logger.Activity, "type", "end", "activityId", activityId)
		// 调用活动结束方法
		handler, err := GetHandler(activityId)
		if err != nil {
			logger.Error(err)
			return
		}
		handler.ActivityEnd()
	}

	// 结束时间在当前时间之后，调度任务
	if endDate.After(now) {
		future, err := a.taskScheduler.Schedule(endTask, endDate)
		if err != nil {
			logger.Error(err)
			return
		}
		// 注册任务
		a.registerTask(activityId, "end", future)
	}
}

// scheduleCronTasks 调度其他特殊时间任务（对应 Java 的 scheduleCronTasks 方法）
func (a *ActivityScheduler) scheduleCronTasks(activity *AbsActivityHandler, activityId int32, now time.Time) {
	// 遍历 Cron 方法映射
	for cronExpr, method := range activity.CronMethods {
		cronDate, _ := schedule.GetNextTriggerTimeAfter(cronExpr, now)
		if !cronDate.After(now) {
			continue
		}

		// 定义 Cron 任务
		cronTask := func() {
			logger.Log(logger.Activity, "type", method.Name, "activityId", activityId)
			// 调用对应方法
			handler, err := GetHandler(activityId)
			if err != nil {
				logger.Error(err)
				return
			}

			// 反射调用方法（对应 Java 的 method.invoke）
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
		future, err := a.taskScheduler.Schedule(cronTask, cronDate)
		if err != nil {
			logger.Error(err)
			continue
		}

		// 注册任务
		a.registerTask(activityId, cronExpr, future)
	}
}

// registerTask 注册任务（对应 Java 的 registerTask 方法）
func (a *ActivityScheduler) registerTask(activityId int32, name string, task Cancellable) {
	// 外层 sync.Map 获取内层任务映射
	innerVal, ok := a.scheduledTasks.Load(activityId)
	var innerMap *innerTaskMap
	if !ok {
		// 不存在则创建新的内层任务映射
		innerMap = &innerTaskMap{
			tasks: make(map[string]Cancellable),
		}
		a.scheduledTasks.Store(activityId, innerMap)
	} else {
		innerMap = innerVal.(*innerTaskMap)
	}

	// 内层加锁更新任务
	innerMap.mu.Lock()
	defer innerMap.mu.Unlock()

	// 如有旧任务，先取消
	if oldTask, exists := innerMap.tasks[name]; exists {
		oldTask.Cancel()
	}
	innerMap.tasks[name] = task
}

// CancelScheduledActivity 取消活动调度（对应 Java 的 cancelScheduledActivity 方法）
func (a *ActivityScheduler) CancelScheduledActivity(activityId int32) {
	// 从外层 sync.Map 移除内层任务映射
	innerVal, ok := a.scheduledTasks.LoadAndDelete(activityId)
	if !ok {
		return
	}

	innerMap := innerVal.(*innerTaskMap)
	innerMap.mu.RLock()
	defer innerMap.mu.RUnlock()

	// 取消所有任务
	for _, task := range innerMap.tasks {
		task.Cancel()
	}
}
