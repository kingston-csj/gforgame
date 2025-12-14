package schedule

import (
	"time"
)

type ScheduleExpressionParser interface {
	// IsValidExpression 校验表达式是否为当前解析器支持的合法格式
	IsValidExpression(expression string) bool

	// GetNextTriggerTimeAfter 计算距离参考时间起的下一次触发时间
	GetNextTriggerTimeAfter(expression string, t time.Time) (time.Time, error)

	// IsPeriodicExpression 是否为周期表达式
	IsPeriodicExpression(expression string) bool
} 