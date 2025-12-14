package schedule

import (
	"errors"
	"io/github/gforgame/examples/system"
	"strconv"
	"strings"
	"time"
)

// ---------------------- 开服天数解析器结构体 ----------------------
// OpenServerScheduleExpressionParser 实现 ScheduleExpressionParser 接口
type OpenServerScheduleExpressionParser struct{}

// IsValidExpression 校验表达式是否合法（固定格式：秒 分 时 天 *，分割后长度为5）
func (o *OpenServerScheduleExpressionParser) IsValidExpression(expression string) bool {
	splits := strings.Split(expression, " ")
	return len(splits) == 5
}

// GetNextTriggerTimeAfter 计算参考时间后，对应开服天数的下一次触发时间
func (o *OpenServerScheduleExpressionParser) GetNextTriggerTimeAfter(expression string, _ time.Time) (time.Time, error) {
	// 获取开服时间
	openServerTime, err := getOpenServerDate()
	if err != nil {
		return time.Time{}, err
	}

	splits := strings.Split(expression, " ")
	if len(splits) != 5 {
		return time.Time{}, nil // 非合法表达式，返回零值时间 
	}

	second, err := strconv.Atoi(splits[0])
	if err != nil {
		return time.Time{}, err
	}
	minute, err := strconv.Atoi(splits[1])
	if err != nil {
		return time.Time{}, err
	}
	hour, err := strconv.Atoi(splits[2])
	if err != nil {
		return time.Time{}, err
	}
	day, err := strconv.Atoi(splits[3])
	if err != nil {
		return time.Time{}, err
	}

	nextTime := openServerTime.
		AddDate(0, 0, day).     
		Add(time.Hour * time.Duration(hour)).   
		Add(time.Minute * time.Duration(minute)).  
		Add(time.Second * time.Duration(second)) 

	return nextTime, nil
}

// getOpenServerDate 获取开服时间（需替换为你的项目实际逻辑）
func getOpenServerDate() (time.Time, error) {
	// 示例：此处返回一个固定时间，实际开发中请替换为从配置/数据库获取开服时间的逻辑
	// 若开服时间不存在，返回 error
	openServerStr := system.GetOpenSeverTime().Data.(string) 
	if openServerStr == "" {
		return time.Time{}, errors.New("open server time is empty")
	}
		
	openServerTime, err := time.Parse("2006-01-02 15:04:05", openServerStr)
	if err != nil {
		return time.Time{}, err
	}
	return openServerTime, nil
}

func (d *OpenServerScheduleExpressionParser) IsPeriodicExpression(expression string) bool {
	return false
}