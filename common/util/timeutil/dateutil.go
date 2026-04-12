package timeutil

import (
	"time"
)

/**
 * 计算指定时间戳与今天相差的天数
 * 如果是今天，返回1， 昨天=2，前天=3...
 * 注意， 如果指定时间戳比当前时间戳大，统一返回-1!
 *
 * @param timestamp 毫秒时间戳
 * @return 今天=1，昨天=2，前天=3... 未来时间返回-1
 */
func GetDayDiffFromToday(timestamp int64) int32 {
	// 获取系统本地时区（和 Java ZoneId.systemDefault() 一致）
	zone := time.Local

	// 今天 0点0分0秒
	today := time.Now().In(zone).Truncate(24 * time.Hour)

	// 把毫秒时间戳转成本地时间
	targetTime := time.UnixMilli(timestamp).In(zone)
	// 目标日期的 0点0分0秒
	targetDate := targetTime.Truncate(24 * time.Hour)

	// 如果是未来时间，返回 -1
	if targetDate.After(today) {
		return -1
	}

	// 计算相差天数
	diffDays := today.Sub(targetDate).Hours() / 24

	// 规则：今天=1，昨天=2...
	return int32(int(diffDays) + 1)
}