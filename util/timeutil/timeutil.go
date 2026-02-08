package timeutil

import (
	"errors"
	"time"
)

// 定义固定格式常量：对应 "yyyy-MM-dd HH:mm:ss"
const (
	// LayoutYmdHms 日期时间格式：年-月-日 时:分:秒（对应示例 "1970-01-01 00:00:00"）
	LayoutYmdHms = "2006-01-02 15:04:05"

	MILLIS_PER_SECOND = int64(1000)
	MILLIS_PER_MINUTE = int64(60) * MILLIS_PER_SECOND
	MILLIS_PER_HOUR   = int64(60) * MILLIS_PER_MINUTE
	MILLIS_PER_DAY    = int64(24) * MILLIS_PER_HOUR
	MILLIS_PER_WEEK   = int64(7) * MILLIS_PER_DAY
)

// ParseLocalTime 将符合 LayoutYmdHms 格式的日期字符串，解析为本地时区（time.Local）的 time.Time
// 参数：
//   dateStr - 待解析的日期字符串，格式必须严格匹配 "yyyy-MM-dd HH:mm:ss"（如 "1970-01-01 00:00:00"）
// 返回：
//   time.Time - 解析后的本地时区时间
//   error - 解析失败错误（格式不匹配、字符串为空、日期无效等）
func ParseLocalTime(dateStr string) (time.Time, error) {
	// 参数校验，避免空字符串传入
	if len(dateStr) == 0 {
		return time.Time{}, errors.New("date string cannot be empty")
	}

	utcTime, err := time.Parse(LayoutYmdHms, dateStr)
	if err != nil {
		return time.Time{}, errors.New("parse date failed: " + err.Error())
	}

	// 将 UTC 时区转换为本地时区（time.Local）
	localTime := utcTime.In(time.Local)

	return localTime, nil
}
 