package activity

import (
	"errors"
	"io/github/gforgame/examples/config"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/logger"
	"io/github/gforgame/schedule"
	"io/github/gforgame/util/timeutil"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)


type Cancellable interface {
	// Cancel 取消任务
	Cancel() bool
}

// 定义错误常量，方便外部判断
var (
	ErrInvalidTask     = errors.New("task scheduler: task function is nil")
	ErrInvalidExecTime = errors.New("task scheduler: execution time is invalid (zero time)")
)

// timerCancellable 是 Cancellable 接口的具体实现，包装 time.Timer 实现取消逻辑
type timerCancellable struct {
	timer    *time.Timer       // 底层定时器
	canceled atomic.Bool       // 标记是否已取消（原子变量保证并发安全）
	done     chan struct{}     // 标记任务是否已执行
}

// Cancel 实现 Cancellable 接口的 Cancel 方法
// 返回值：true=成功取消（任务未执行），false=任务已执行/已取消/定时器已触发
func (tc *timerCancellable) Cancel() bool {
	// 双重检查+原子操作，避免重复取消
	if tc.canceled.CompareAndSwap(false, true) {
		// 停止定时器：Stop() 返回true表示定时器未触发，false表示已触发/已停止
		stopped := tc.timer.Stop()
		// 关闭一个无缓冲通道后，所有监听该通道的 case <-chan 会立即触发（即使通道中没有数据）
		close(tc.done) 
		return stopped
	}
	return false
}

// TaskScheduler 定时任务调度接口
type TaskScheduler interface {
	// Schedule 调度一次性任务，在指定时间执行
	Schedule(task func(), execTime time.Time) (Cancellable, error)
}

// newTimerCancellable 创建 timerCancellable 实例
func newTimerCancellable(timer *time.Timer) *timerCancellable {
	return &timerCancellable{
		timer: timer,
		done:  make(chan struct{}),
	}
}

// DefaultTaskScheduler 是 TaskScheduler 接口的默认实现
// 基于 time.Timer 实现一次性定时任务调度，支持任务取消、并发安全
type DefaultTaskScheduler struct{}

// NewDefaultTaskScheduler 创建 DefaultTaskScheduler 实例
func NewDefaultTaskScheduler() *DefaultTaskScheduler {
	return &DefaultTaskScheduler{}
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

// Schedule 实现 TaskScheduler 接口的 Schedule 方法
// 功能：调度一次性任务，在指定时间 execTime 执行
// 参数：
//   - task: 要执行的任务函数（不能为 nil）
//   - execTime: 任务执行时间（不能为零值）
// 返回值：
//   - Cancellable: 任务取消器，可调用 Cancel() 取消未执行的任务
//   - error: 入参非法时返回非 nil 错误
func (d *DefaultTaskScheduler) Schedule(task func(), execTime time.Time) (Cancellable, error) {
	// 入参校验
	if task == nil {
		return nil, ErrInvalidTask
	}
	if execTime.IsZero() {
		return nil, ErrInvalidExecTime
	}

	// 计算当前时间到执行时间的延迟
	now := time.Now()
	delay := execTime.Sub(now)

	// 处理延迟<=0的情况（执行时间已过期，立即异步执行任务）
	if delay <= 0 {
		logger.Info("task scheduler: execTime is in the past, execute task immediately")
		// 异步执行任务，避免阻塞当前goroutine
		go func() {
			defer func() {
				// 捕获任务执行中的panic，避免程序崩溃
				if r := recover(); r != nil {
					logger.Error2("task scheduler: task panic" , r.(error))
				}
			}()
			task()
		}()

		// 返回一个已取消的Cancellable（任务已执行，无法取消）
		return &timerCancellable{
			canceled: atomic.Bool{},
		}, nil
	}

	// 延迟执行：创建定时器，调度任务
	timer := time.NewTimer(delay)
	cancellable := newTimerCancellable(timer)

	// 启动goroutine等待定时器触发，执行任务
	go func() {
		defer func() {
			// 捕获任务执行中的panic
			if r := recover(); r != nil {
				logger.Error2("task scheduler: task panic", r.(error))
			}
		}()

		select {
		case <-timer.C:
			// 定时器触发，执行任务
			task()
			// 标记任务已执行
			cancellable.canceled.Store(true)
			close(cancellable.done)
		case <-cancellable.done:
			// 任务被取消，直接退出
			return
		}
	}()

	return cancellable, nil
}

// NewActivityScheduler 创建新的活动调度器
func NewActivityScheduler(taskScheduler TaskScheduler) *ActivityScheduler {
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
	referDate := a.getReferenceDate(getter, startCron)
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
	getter.GetBaseHandler().EndTime = endDate.UnixMilli()

	// 如果结束时间在当前时间之后，调度相关任务
	if endDate.After(now) {
		a.scheduleStartTask(getter, startDate, now, endDate)
		a.scheduleEndTask(activityId, endDate, now)
		a.scheduleCronTasks(getter, activityId, now)
	}
}

// getReferenceDate 获取参考日期
func (a *ActivityScheduler) getReferenceDate(handler BaseHandlerGetter, startCron string) time.Time {
	// 程序运行期间，参考日期以最近的活动结束时间为准
	if handler.GetBaseHandler().LastEndDate != nil {
		return *handler.GetBaseHandler().LastEndDate
	}

	// 周期性活动，取当前时间
	if schedule.GetParser(startCron).IsPeriodicExpression(startCron) {
		return time.Now()
	}

	// 刚启服场景，取 1970-01-01 作为参考日期
	defaultDate, _ := timeutil.ParseLocalTime("1970-01-01 00:00:00")
	return defaultDate
}

// scheduleStartTask 调度活动开始任务
func (a *ActivityScheduler) scheduleStartTask(
	handler BaseHandlerGetter,
	startDate, now, endDate time.Time) {
	// 设置活动结束时间
	handler.GetBaseHandler().EndTime = endDate.UnixMilli()
	activityId := handler.GetBaseHandler().ActivityId	
	// 定义开始任务
	startTask := func() {
		logger.Log(logger.Activity, "type", "start", "activityId", activityId)
		// 调用活动开始方法
		handler.GetBaseHandler().ActivityStart()
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
func (a *ActivityScheduler) scheduleCronTasks(handler BaseHandlerGetter, activityId int32, now time.Time) {
	// 遍历 Cron 方法映射
	for cronExpr, method := range handler.GetBaseHandler().CronMethods {
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
