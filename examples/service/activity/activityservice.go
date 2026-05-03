package activity

import (
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/protos"
	"sync"
)

type ActivityService struct {
	activityScheduler *ActivityScheduler
}

var (
	instance          *ActivityService
	once              sync.Once
	activityScheduler *ActivityScheduler = NewActivityScheduler(context.TaskScheduler)
)

func GetActivityService() *ActivityService {
	once.Do(func() {
		instance = &ActivityService{
			activityScheduler: activityScheduler,
		}

		firstRechargeHandler := NewFirstRechargeActivityHandler(activityScheduler)
		registerHandler(1001, firstRechargeHandler)
		registerHandler(1002, firstRechargeHandler)

	})
	return instance
}

func (s *ActivityService) ScheduleAllActivity() {
	for _, activityData := range config.QueryAll[configdomain.ActivityData]() {
		handler, err := GetHandler(activityData.Id)
		if err != nil {
			panic(err)
		}
		s.activityScheduler.ScheduleActivity(handler)
	}
}

func (sv *ActivityService) OnPlayerLogin(p *playerdomain.Player) {
	activityVos := make([]*protos.ActivityVo, 0)
	for _, activityData := range config.QueryAll[configdomain.ActivityData]() {
		handler, err := GetHandler(activityData.Id)
		if err != nil {
			panic(err)
		}
		if handler.IsInActivity(p) {
			activityVo := handler.LoadActivityInfo(p)
			activityVos = append(activityVos, activityVo)
		}

	}
	push := &protos.PushActivityLoadAll{
		ActivityVos: activityVos,
	}
	io.NotifyPlayer(p, push)
}

func GetActivityRewards(activityId int32) []*configdomain.ActivityRewardData {
	rewardDatas := make([]*configdomain.ActivityRewardData, 0)
	for _, rewardData := range config.QueryAll[configdomain.ActivityRewardData]() {
		if rewardData.ActivityId == activityId {
			rewardDatas = append(rewardDatas, rewardData)
		}
	}
	return rewardDatas
}
