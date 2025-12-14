package activity

import (
	"sync"

	"io/github/gforgame/data"
	"io/github/gforgame/examples/config"
	configdomain "io/github/gforgame/examples/domain/config"
)

type ActivityService struct {
	activityScheduler *ActivityScheduler
}

var (
	instance *ActivityService
	once     sync.Once
	taskScheduler TaskScheduler = NewDefaultTaskScheduler()
	activityScheduler *ActivityScheduler = NewActivityScheduler(taskScheduler)
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
	container := config.QueryContainer[configdomain.ActivityData,*data.Container[int32, configdomain.ActivityData]]()
	for _, activityData := range container.GetAllRecords() {
		handler, err := GetHandler(activityData.Id)
		if err != nil {
			panic(err)
		}
		s.activityScheduler.ScheduleActivity(handler)
	}
}
