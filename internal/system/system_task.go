package system

import (
	"fmt"
	"log"
	"time"

	"github.com/forfun/gforgame/common/eventbus"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/timeutil"
	"github.com/forfun/gforgame/internal/events"

	"github.com/robfig/cron/v3"
)

// StartSystemTask 启动系统定时任务
func (s *SystemService) StartSystemTask() {
	// 创建一个新的cron调度器，使用本地时区
	s.scheduler = cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))

	// 添加每日23:59:59执行任务
	_, err := s.scheduler.AddFunc("59 59 23 * * *", s.performDailyUpdate)
	if err != nil {
		log.Printf("添加每日更新任务失败: %v", err)
	}

	// 每小时执行的任务示例
	_, err = s.scheduler.AddFunc("0 0 * * * *", s.performHourlyUpdate)
	if err != nil {
		log.Printf("添加每小时更新任务失败: %v", err)
	}

	// 每周日晚?11:59:59 执行的定时任任务
	_, err = s.scheduler.AddFunc("59 59 23 * * 0", s.performWeeklyUpdate)
	if err != nil {
		log.Printf("添加每周更新任务失败: %v", err)
	}

	// 每个月最后一天的 23:59:59 执行的定时任务（每个月最后几天定时触发，二次判断）
	_, err = s.scheduler.AddFunc("59 59 23 28-31 * *", s.performMonthlyUpdate)
	if err != nil {
		log.Printf("添加每月更新任务失败: %v", err)
	}

	// 启动调度器
	s.scheduler.Start()

	logger.Info("系统定时任务已启动")
}

// StopSystemTask 停止系统定时任务
func (s *SystemService) StopSystemTask() {
	if s.scheduler != nil {
		s.scheduler.Stop()
		log.Println("系统定时任务已停止")
	}
}

// performDailyUpdate 执行每日更新操作
func (s *SystemService) performDailyUpdate() {
	log.Println("执行每日更新任务 -", time.Now().Format(timeutil.LayoutYmdHms))

	dailyReset := s.GetDailyReset()

	// 更新每日重置时间
	newResetTimestamp := int64(time.Now().Unix())
	dailyReset.Save(newResetTimestamp)
	fmt.Printf("当前每日重置时间: %d\n", newResetTimestamp)
	eventbus.Default().Publish(events.SystemDailyReset, newResetTimestamp)
}

// PerformDailyUpdate 提供给外部主动触发每日重置的入口。
func (s *SystemService) PerformDailyUpdate() {
	s.performDailyUpdate()
}

// performHourlyUpdate 执行每小时更新操作
func (s *SystemService) performHourlyUpdate() {
	log.Println("执行每小时更新任务 -", time.Now().Format(timeutil.LayoutYmdHms))

	// 在这里添加需要每小时更新的逻辑
	// 例如：更新在线玩家状态、检查服务器负载?
}

// performWeeklyUpdate 执行每周更新操作
func (s *SystemService) performWeeklyUpdate() {
	log.Println("执行每周更新任务 -", time.Now().Format(timeutil.LayoutYmdHms))

	weeklyReset := s.GetWeeklyReset()

	// 更新每周重置时间
	newResetTimestamp := int64(time.Now().Unix())
	weeklyReset.Save(newResetTimestamp)
	fmt.Printf("当前每周重置时间: %d\n", newResetTimestamp)
	eventbus.Default().Publish(events.SystemWeeklyReset, newResetTimestamp)
}

// performMonthlyUpdate 执行每月更新操作
func (s *SystemService) performMonthlyUpdate() {
	// cron库不支持每个月最后一天的表达式，这里加一个二次判断
	// 如果当前日期不是1号，说明不是每个月的最后一天，直接返回
	now := time.Now()
	if now.AddDate(0, 0, 1).Day() != 1 {
		return
	}
	log.Println("执行每月更新任务 -", time.Now().Format(timeutil.LayoutYmdHms))

	monthlyReset := s.GetMonthlyReset()

	// 更新每月重置时间
	newResetTimestamp := int64(time.Now().Unix())
	monthlyReset.Save(newResetTimestamp)
	fmt.Printf("当前每月重置时间? %d\n", newResetTimestamp)
	eventbus.Default().Publish(events.SystemMonthlyReset, newResetTimestamp)
}

// AddCustomTask 添加自定义定时任务
// cron表达式格式
// 例如: "0 30 9 * * *" 表示每天上午9:30执行
func (s *SystemService) AddCustomTask(spec string, task func()) error {
	if s.scheduler == nil {
		s.scheduler = cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))
		s.scheduler.Start()
	}

	_, err := s.scheduler.AddFunc(spec, task)
	if err != nil {
		log.Printf("添加自定义任务失败: %v", err)
		return err
	}

	log.Printf("成功添加自定义任务，执行计划: %s", spec)
	return nil
}
