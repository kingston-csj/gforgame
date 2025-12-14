package system

var (
	dailyReset   *DailyReset
	weeklyReset  *WeeklyReset
	monthlyReset *MonthlyReset
	openSever    *OpenSeverTime
)

func init() {
	once.Do(func() {
		GetSystemParameterService().init()

		// 从数据库加载数据
		loadParameterData("1001")
		loadParameterData("1002")
		loadParameterData("1003")
		loadParameterData("1004")
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
