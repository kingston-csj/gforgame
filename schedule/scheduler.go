package schedule

import (
	"errors"
	"fmt"
	"io/github/gforgame/logger"
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
	// Schedule 调度一次性任务，在指定延迟（毫秒）后执行
	// 参数：
	//   - task: 要执行的任务函数（不能为 nil）
	//   - delay: 任务执行延迟时间（毫秒），必须大于0
	// 返回值：
	//   - Cancellable: 任务取消器，可调用 Cancel() 取消未执行的任务
	//   - error: 入参非法时返回非 nil 错误
	Schedule(task func(), delay int64) (Cancellable, error)
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

// Schedule 实现 TaskScheduler 接口的 Schedule 方法
func (d *DefaultTaskScheduler) Schedule(task func(), delay int64) (Cancellable, error) {
	// 入参校验
	if task == nil {
		return nil, ErrInvalidTask
	}

	// 处理延迟<=0的情况（执行时间已过期，立即异步执行任务）
	if delay <= 0 {
		logger.Info("task scheduler: delay <= 0, execute task immediately")
		// 异步执行任务，避免阻塞当前goroutine
		go func() {
			defer func() {
				// 捕获任务执行中的panic，避免程序崩溃
				if r := recover(); r != nil {
					logger.Error(fmt.Errorf("task panic: %v", r))
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
	timer := time.NewTimer(time.Duration(delay) * time.Millisecond)
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