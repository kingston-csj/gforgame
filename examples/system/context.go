package system

var (
	dailyReset *DailyReset
)

func init() {
	once.Do(func() {
		dailyReset = &DailyReset{
			ID: "1001",
		}
		// 从数据库加载数据
		loadParameterData(dailyReset.GetID())

	})
}

func loadParameterData(param string) {
	GetSystemParameterService().GetOrCreateSystemParameterRecord(param)
}

func GetDailyReset() *DailyReset {
	return dailyReset
}
