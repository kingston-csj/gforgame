package schedule

import (
	"time"
)

var (
	parsers []ScheduleExpressionParser 
)

func AddParserBefore(parser ScheduleExpressionParser) {
	parsers = append([]ScheduleExpressionParser{parser}, parsers...)
}

func AddParserAfter(parser ScheduleExpressionParser) {
	parsers = append(parsers, parser)
}

func init() {
	AddParserBefore(NewCronParser())
}

// GetNextTriggerTimeAfter 按照解析链,逐一解析表达式,如果表达式符合规则,则按当前节点解析器进行解析
func GetNextTriggerTimeAfter(expression string, t time.Time) (time.Time, error) {
	for _, parser := range parsers {
		if !parser.IsValidExpression(expression) {
			continue
		}
		nextTime, err := parser.GetNextTriggerTimeAfter(expression, t)
		if err != nil {
			return time.Time{}, err
		}
		if nextTime.After(t) {
			return nextTime, nil
		}
	}
	return time.Time{}, nil
}

func GetParser(expression string) ScheduleExpressionParser {
	for _, parser := range parsers {
		if parser.IsValidExpression(expression) {
			return parser
		}
	}
	return nil
}
