package system

import (
	"fmt"
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/events"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	// 全局cron调度器
	scheduler *cron.Cron
)

// StartSystemTask 启动系统定时任务
func StartSystemTask() {
	// 创建一个新的cron调度器，使用本地时区
	scheduler = cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))

	// 添加每日23:59:59执行的任务
	// 秒 分 时 日 月 星期
	_, err := scheduler.AddFunc("59 59 23 * * *", performDailyUpdate)
	if err != nil {
		log.Printf("添加每日更新任务失败: %v", err)
	}

	// 每小时执行的任务示例
	_, err = scheduler.AddFunc("0 0 * * * *", performHourlyUpdate)
	if err != nil {
		log.Printf("添加每小时更新任务失败: %v", err)
	}

	// 每周日晚上 11:59:59 执行的定时任务
	_, err = scheduler.AddFunc("59 59 23 * * 0", performWeeklyUpdate)
	if err != nil {
		log.Printf("添加每周更新任务失败: %v", err)
	}

	// 每个月最后一天的 23:59:59 执行的定时任务(暂时无法工作)
	_, err = scheduler.AddFunc("59 59 23 L * *", performMonthlyUpdate)
	if err != nil {
		log.Printf("添加每月更新任务失败: %v", err)
	}

	// 启动调度器
	scheduler.Start()

	log.Println("系统定时任务已启动")
}

// StopSystemTask 停止系统定时任务
func StopSystemTask() {
	if scheduler != nil {
		scheduler.Stop()
		log.Println("系统定时任务已停止")
	}
}

// performDailyUpdate 执行每日更新操作
func performDailyUpdate() {
	log.Println("执行每日更新任务 -", time.Now().Format("2006-01-02 15:04:05"))

	dailyReset := GetDailyReset()

	// 更新每日重置时间戳
	newResetTimestamp := int64(time.Now().Unix())
	dailyReset.Save(newResetTimestamp)
	fmt.Printf("当前每日重置时间戳: %d\n", newResetTimestamp)
	context.EventBus.Publish(events.SystemDailyReset, newResetTimestamp)
}

// performHourlyUpdate 执行每小时更新操作
func performHourlyUpdate() {
	log.Println("执行每小时更新任务 -", time.Now().Format("2006-01-02 15:04:05"))

	// 在这里添加需要每小时更新的逻辑
	// 例如：更新在线玩家状态、检查服务器负载等
}

// performWeeklyUpdate 执行每周更新操作
func performWeeklyUpdate() {
	log.Println("执行每周更新任务 -", time.Now().Format("2006-01-02 15:04:05"))

	weeklyReset := GetWeeklyReset()

	// 更新每周重置时间戳
	newResetTimestamp := int64(time.Now().Unix())
	weeklyReset.Save(newResetTimestamp)
	fmt.Printf("当前每周重置时间戳: %d\n", newResetTimestamp)
	context.EventBus.Publish(events.SystemWeeklyReset, newResetTimestamp)
}

// performMonthlyUpdate 执行每月更新操作
func performMonthlyUpdate() {
	log.Println("执行每月更新任务 -", time.Now().Format("2006-01-02 15:04:05"))

	monthlyReset := GetMonthlyReset()

	// 更新每月重置时间戳
	newResetTimestamp := int64(time.Now().Unix())
	monthlyReset.Save(newResetTimestamp)
	fmt.Printf("当前每月重置时间戳: %d\n", newResetTimestamp)
	context.EventBus.Publish(events.SystemMonthlyReset, newResetTimestamp)
}

// AddCustomTask 添加自定义定时任务
// cron表达式格式: 秒 分 时 日 月 星期
// 例如: "0 30 9 * * *" 表示每天上午9:30执行
func AddCustomTask(spec string, task func()) error {
	if scheduler == nil {
		scheduler = cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))
		scheduler.Start()
	}

	_, err := scheduler.AddFunc(spec, task)
	if err != nil {
		log.Printf("添加自定义任务失败: %v", err)
		return err
	}

	log.Printf("成功添加自定义任务，执行计划: %s", spec)
	return nil
}
