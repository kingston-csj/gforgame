package system

import (
	"github.com/forfun/gforgame/common/schedule"
	systemrepo "github.com/forfun/gforgame/internal/infra/repository/system"
	"github.com/robfig/cron/v3"
)

type SystemService struct {
	dailyReset   *DailyReset
	weeklyReset  *WeeklyReset
	monthlyReset *MonthlyReset
	openSever    *OpenSeverTime
	scheduler    *cron.Cron
	repository   *systemrepo.SystemRepository
}

func NewSystemService(repo *systemrepo.SystemRepository) *SystemService {
	s := &SystemService{}
	s.repository = repo
	s.dailyReset = NewDailyReset(repo)
	s.weeklyReset = NewWeeklyReset(repo)
	s.monthlyReset = NewMonthlyReset(repo)
	s.openSever = NewOpenServerTime(repo)
	return s
}

// Init 初始化系统参数依赖与调度表达式解析器。
func (s *SystemService) Init() {
	schedule.AddParserAfter(NewOpenServerScheduleExpressionParser(s))
}

func (s *SystemService) GetDailyReset() *DailyReset {
	return s.dailyReset
}

func (s *SystemService) GetWeeklyReset() *WeeklyReset {
	return s.weeklyReset
}

func (s *SystemService) GetMonthlyReset() *MonthlyReset {
	return s.monthlyReset
}

func (s *SystemService) GetOpenSeverTime() *OpenSeverTime {
	return s.openSever
}
