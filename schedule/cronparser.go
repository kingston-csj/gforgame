package schedule

import (
	"time"

	"github.com/robfig/cron/v3"
)
type CronParser struct {
	cronParser *cron.Parser
}

func NewCronParser() *CronParser {
	return &CronParser{}
}
 

// IsValidExpression 实现接口的校验方法
func (d *CronParser) IsValidExpression(expression string) bool {
	// 尝试解析表达式，无错误则为合法格式
	_, err := d.cronParser.Parse(expression)
	return err == nil
}

// GetNextTriggerTimeAfter 实现接口的下一次触发时间计算方法
func (d *CronParser) GetNextTriggerTimeAfter(expression string, t time.Time) (time.Time, error) {
	// 解析 cron 表达式
	schedule, err := d.cronParser.Parse(expression)
	if err != nil {
		return time.Time{}, err // 解析失败，返回零值时间和错误
	}

	// 调用 cron.Schedule 的 Next 方法，获取参考时间后的下一次触发时间
	nextTriggerTime := schedule.Next(t)
	return nextTriggerTime, nil
}

func (d *CronParser) IsPeriodicExpression(expression string) bool {
	// 解析Cron表达式
	schedule, err := d.cronParser.Parse(expression)
	if err != nil {
		panic("解析Cron表达式失败: " + err.Error())
	}

	// 说明：robfig/cron 的 Schedule 解析时，若未指定时区，默认使用本地时区（time.Local）
	// 若需严格指定默认时区，可在创建 cron.Parser 时通过 opts 配置 
	now := time.Now()
	// 获取第一个触发时间
	firstFireTime := schedule.Next(now)
	// 判断第一个触发时间是否存在
	if firstFireTime.IsZero() {
		return false
	}

	// 获取第二个触发时间
	secondFireTime := schedule.Next(firstFireTime)
	//存在第二个触发时间则为周期性
	return !secondFireTime.IsZero()
}