package system

import (
	"github.com/forfun/gforgame/common/schedule"
)

var (
	dailyReset   *DailyReset
	weeklyReset  *WeeklyReset
	monthlyReset *MonthlyReset
	openSever    *OpenSeverTime
)

func init() {
	once.Do(func() {
		// 从数据库加载数据
		GetSystemParameterService().init()
		dailyReset = NewDailyReset()
		weeklyReset = NewWeeklyReset()
		monthlyReset = NewMonthlyReset()
		openSever = NewOpenServerTime()

		schedule.AddParserAfter(&OpenServerScheduleExpressionParser{})
	})
}

func loadParameterData(param string) {
	GetSystemParameterService().GetOrCreateSystemParameterRecord(param)
}

func GetDailyReset() *DailyReset {
	return dailyReset
}

func GetWeeklyReset() *WeeklyReset {
	return weeklyReset
}

func GetMonthlyReset() *MonthlyReset {
	return monthlyReset
}

func GetOpenSeverTime() *OpenSeverTime {
	return openSever
}
