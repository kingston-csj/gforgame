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
		dailyReset = &DailyReset{}
		weeklyReset = &WeeklyReset{}
		monthlyReset = &MonthlyReset{}
		openSever = &OpenSeverTime{}
		// 从数据库加载数据
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
