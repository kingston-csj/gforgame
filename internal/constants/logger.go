package constants

type LoggerCategory struct {
	Code    string
	LogType int
	Desc    string
}

const (
	LogTypeOperation = 0
	LogTypeSystem    = 1
)

var (
	LoggerAdmin = LoggerCategory{
		Code:    "admin",
		LogType: LogTypeSystem,
		Desc:    "后台管理",
	}

	LoggerApplication = LoggerCategory{
		Code:    "application",
		LogType: LogTypeSystem,
		Desc:    "应用系统",
	}

	LoggerPlayer = LoggerCategory{
		Code:    "player",
		LogType: LogTypeOperation,
		Desc:    "玩家",
	}

	LoggerNet = LoggerCategory{
		Code:    "net",
		LogType: LogTypeOperation,
		Desc:    "网络",
	}

	LoggerMail = LoggerCategory{
		Code:    "mail",
		LogType: LogTypeOperation,
		Desc:    "邮件",
	}

	LoggerItem = LoggerCategory{
		Code:    "item",
		LogType: LogTypeOperation,
		Desc:    "道具",
	}

	LoggerHero = LoggerCategory{
		Code:    "hero",
		LogType: LogTypeOperation,
		Desc:    "英雄",
	}

	LoggerActivity = LoggerCategory{
		Code:    "activity",
		LogType: LogTypeSystem,
		Desc:    "活动调度",
	}

	LoggerDebug = LoggerCategory{
		Code:    "debug",
		LogType: LogTypeSystem,
		Desc:    "调试",
	}
)

// loggerCategoriesAll 枚举所有已定义的日志分类（仿 Kotlin enum entries 遍历）。
// 新增 LoggerCategory 时记得在下面的切片里追加一行。
func loggerCategoriesAll() []LoggerCategory {
	return []LoggerCategory{
		LoggerAdmin,
		LoggerApplication,
		LoggerPlayer,
		LoggerMail,
		LoggerItem,
		LoggerHero,
		LoggerActivity,
		LoggerDebug,
	}
}

// OperationLogCodes 返回所有运营日志（logType == 0）的小写 Code 列表。
// 对应 Kotlin 接口：
//
//	LoggerFunction.entries.filter { it.logType == 0 }
//	                    .map { it.name.lowercase() }
func OperationLogCodes() []string {
	all := loggerCategoriesAll()
	result := make([]string, 0, len(all))
	for _, cat := range all {
		if cat.LogType == LogTypeOperation {
			result = append(result, cat.Code)
		}
	}
	return result
}
