package activity

import (
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	"sync"
)

type ActivityService struct {
	activityScheduler *ActivityScheduler
}

var (
	instance *ActivityService
	once     sync.Once
	activityScheduler *ActivityScheduler = NewActivityScheduler(context.TaskScheduler)
)

func GetActivityService() *ActivityService {
	once.Do(func() {
		instance = &ActivityService{
			activityScheduler: activityScheduler,
		}

		firstRechargeHandler := NewFirstRechargeActivityHandler(activityScheduler);
		registerHandler(1001, firstRechargeHandler)

	})
	return instance
}

func (s *ActivityService) ScheduleAllActivity( ) {
	for _, activityData := range config.QueryAll[configdomain.ActivityData]() {
		handler, err := GetHandler(activityData.Id)
		if err != nil {
			panic(err)
		}
		s.activityScheduler.ScheduleActivity(handler)
	}
}
