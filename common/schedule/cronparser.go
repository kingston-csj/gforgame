package schedule

import (
	"fmt"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
)

type CronParser struct {
	cronParser *cron.Parser
}

func NewCronParser() *CronParser {
	parser := cron.NewParser(
		cron.Second | // 解析秒
		cron.Minute | // 解析分
		cron.Hour |   // 解析时
		cron.Dom |    // 解析日
		cron.Month |  // 解析月
		cron.Dow,     // 解析周（无Year选项）
	)
	return &CronParser{
		cronParser: &parser,
	}
}

// ---------------- 新增辅助函数：处理Quartz 7字段表达式 ----------------
// splitCronFields 拆分cron表达式为字段切片（按空格分割）
func splitCronFields(expr string) []string {
	var fields []string
	field := ""
	for _, c := range expr {
		if c == ' ' {
			if field != "" {
				fields = append(fields, field)
				field = ""
			}
		} else {
			field += string(c)
		}
	}
	if field != "" {
		fields = append(fields, field)
	}
	return fields
}

// replaceQuestionMark 将Quartz的?替换为Go cron的*（仅用于Dow/Dom字段）
func replaceQuestionMark(field string) string {
	if field == "?" {
		return "*"
	}
	return field
}

// parseYear 解析年字段，返回合法的年份（0表示不限制年份）
func parseYear(yearStr string) (int, error) {
	if yearStr == "*" {
		return 0, nil // * 表示不限制年份
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, fmt.Errorf("invalid year: %s, err: %v", yearStr, err)
	}
	if year < 1970 || year > 9999 { // 限制年份范围，避免无效值
		return 0, fmt.Errorf("year %d out of range [1970, 9999]", year)
	}
	return year, nil
}

// get6FieldCron 从7字段Quartz表达式提取6字段Go cron表达式（替换?）
func get6FieldCron(fields []string) string {
	return fmt.Sprintf("%s %s %s %s %s %s",
		fields[0],
		fields[1],
		fields[2],
		fields[3],
		fields[4],
		replaceQuestionMark(fields[5]), // Dow字段替换?为*
	)
}

// getNextValidTimeWithYear 计算符合目标年份的下一次触发时间
func (d *CronParser) getNextValidTimeWithYear(schedule cron.Schedule, targetYear int, t time.Time) (time.Time, error) {
	if targetYear == 0 { // 不限制年份，直接返回cron计算的时间
		return schedule.Next(t), nil
	}

	nextTime := schedule.Next(t)
	// 循环查找目标年份的时间，最多循环100次避免死循环
	maxLoop := 100
	loopCount := 0
	for nextTime.Year() != targetYear {
		loopCount++
		if loopCount > maxLoop {
			return time.Time{}, fmt.Errorf("exceed max loop count (100), no valid time in year %d", targetYear)
		}

		// 目标年份已过，无有效时间
		if nextTime.Year() > targetYear {
			return time.Time{}, fmt.Errorf("no valid time in year %d (next time is %s)", targetYear, nextTime.Format("2006-01-02"))
		}

		// 跳到目标年份的1月1日，继续计算
		nextTime = time.Date(targetYear, 1, 1, 0, 0, 0, 0, t.Location())
		nextTime = schedule.Next(nextTime)

		// 若跳到目标年后，下一次时间仍大于目标年，说明无有效时间
		if nextTime.Year() > targetYear {
			return time.Time{}, fmt.Errorf("no valid time in year %d", targetYear)
		}
	}

	return nextTime, nil
}

// ---------------- 重写接口方法：适配7字段表达式 ----------------
// IsValidExpression 实现接口的校验方法（兼容6/7字段）
func (d *CronParser) IsValidExpression(expression string) bool {
	fields := splitCronFields(expression)

	// 情况1：7字段（Quartz风格：秒 分 时 日 月 周 年）
	if len(fields) == 7 {
		// 校验年字段合法性
		if _, err := parseYear(fields[6]); err != nil {
			return false
		}
		// 提取6字段cron并校验
		FieldExpr := get6FieldCron(fields)
		_, err := d.cronParser.Parse(FieldExpr)
		return err == nil
	}

	// 情况2：6字段（标准Go cron）
	if len(fields) == 6 {
		_, err := d.cronParser.Parse(expression)
		return err == nil
	}

	// 其他字段长度：非法
	return false
}

// GetNextTriggerTimeAfter 实现接口的下一次触发时间计算方法（兼容6/7字段）
func (d *CronParser) GetNextTriggerTimeAfter(expression string, t time.Time) (time.Time, error) {
	fields := splitCronFields(expression)

	// 情况1：7字段（Quartz风格，含年）
	if len(fields) == 7 {
		// 解析年字段
		targetYear, err := parseYear(fields[6])
		if err != nil {
			return time.Time{}, fmt.Errorf("parse year failed: %v", err)
		}
		// 提取6字段cron并解析
		FieldExpr := get6FieldCron(fields)
		schedule, err := d.cronParser.Parse(FieldExpr)
		if err != nil {
			return time.Time{}, fmt.Errorf("parse 6-field cron failed: %v", err)
		}
		// 计算符合年份的下一次时间
		return d.getNextValidTimeWithYear(schedule, targetYear, t)
	}

	// 情况2：6字段（标准Go cron）
	if len(fields) == 6 {
		schedule, err := d.cronParser.Parse(expression)
		if err != nil {
			return time.Time{}, fmt.Errorf("parse cron failed: %v", err)
		}
		return schedule.Next(t), nil
	}

	// 其他字段长度：非法
	return time.Time{}, fmt.Errorf("invalid cron expression length: %d (expect 6 or 7)", len(fields))
}

// IsPeriodicExpression 判断是否为周期表达式（兼容6/7字段）
func (d *CronParser) IsPeriodicExpression(expression string) bool {
	fields := splitCronFields(expression)

	// 情况1：7字段（Quartz风格，含年）
	if len(fields) == 7 {
		// 解析年字段
		targetYear, err := parseYear(fields[6])
		if err != nil {
			return false // 解析失败，视为非周期
		}
		// 提取6字段cron并解析
		FieldExpr := get6FieldCron(fields)
		schedule, err := d.cronParser.Parse(FieldExpr)
		if err != nil {
			return false
		}

		now := time.Now()
		// 获取第一个触发时间
		firstFireTime, err := d.getNextValidTimeWithYear(schedule, targetYear, now)
		if err != nil || firstFireTime.IsZero() {
			return false
		}

		// 获取第二个触发时间（需基于第一个时间，且仍在目标年）
		secondFireTime, err := d.getNextValidTimeWithYear(schedule, targetYear, firstFireTime)
		if err != nil || secondFireTime.IsZero() {
			return false
		}

		// 存在第二个触发时间 → 周期性
		return true
	}

	// 情况2：6字段（标准Go cron）
	if len(fields) == 6 {
		schedule, err := d.cronParser.Parse(expression)
		if err != nil {
			return false // 解析失败，视为非周期（移除原panic，避免程序崩溃）
		}

		now := time.Now()
		firstFireTime := schedule.Next(now)
		if firstFireTime.IsZero() {
			return false
		}

		secondFireTime := schedule.Next(firstFireTime)
		return !secondFireTime.IsZero()
	}

	// 其他字段长度：非周期
	return false
}