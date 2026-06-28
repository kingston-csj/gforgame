package activity

import (
	"github.com/forfun/gforgame/internal/config"
	"github.com/forfun/gforgame/internal/context"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/protos"
)


type ActivityService struct {
	activityScheduler *ActivityScheduler
}

var activityScheduler = NewActivityScheduler(context.TaskScheduler)

func NewActivityService() *ActivityService {
	s := &ActivityService{
		activityScheduler: activityScheduler,
	}
	firstRechargeHandler := NewFirstRechargeActivityHandler(activityScheduler)
	registerHandler(1001, firstRechargeHandler)
	return s
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
