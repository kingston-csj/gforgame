package activity

import (
	"github.com/forfun/gforgame/common/util/conv"
	"github.com/forfun/gforgame/common/util/timeutil"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/protos"
)

type FirstRechargeActivityHandler struct {
	*BaseActivityHandler
}

func (d *FirstRechargeActivityHandler) GetBaseHandler() *BaseActivityHandler {
	return d.BaseActivityHandler
}

func (d *FirstRechargeActivityHandler) LoadActivityInfo(p *playerdomain.Player) *protos.ActivityVo {
	activityInfo := p.ActivityBox.Data[d.ActivityId]
	activityRewards := GetActivityRewards(d.ActivityId)
	
	if activityInfo == nil {
		activityInfo = &playerdomain.ActivityInfo{
			Rewards: make(map[int32]string),
		}
		p.ActivityBox.Data[d.ActivityId] = activityInfo
		for _, rewardData := range activityRewards {
			activityInfo.Rewards[rewardData.Id] = "0"
		}
	} else {
		firstReward := p.RechargeBox.FirstRechargeTime>0
		if firstReward {
			dayDiff := timeutil.GetDayDiffFromToday(p.RechargeBox.FirstRechargeTime)
			for _, rewardData := range activityRewards {
				if dayDiff >= conv.Int32Value(rewardData.Condition) {
					if activityInfo.Rewards[rewardData.Id] != "2" {
						activityInfo.Rewards[rewardData.Id] = "1"
					}
				}
			}
		}
	}
	activityVo := &protos.ActivityVo{
		ActivityId: d.ActivityId,
		RewardVos: make([]*protos.ActivityRewardVo, 0),
	}
	for rewardId, rewardStatus := range activityInfo.Rewards {
		activityRewardVo := &protos.ActivityRewardVo{
			Id: rewardId,
			Value: rewardStatus,
		}
		activityVo.RewardVos = append(activityVo.RewardVos, activityRewardVo)
	}
	return activityVo
}



func NewFirstRechargeActivityHandler(sched *ActivityScheduler) *FirstRechargeActivityHandler {
	baseHandler := &BaseActivityHandler{
		ActivitySched: sched,
		ActivityId:    1001,
		OnActivityStart: func() error {
			// 首次充值活动的专属启动逻辑
			return nil
		},
		OnActivityEnd: func() error {
			// 首次充值活动的专属结束逻辑
			return nil
		},
	}

	return &FirstRechargeActivityHandler{
		BaseActivityHandler: baseHandler,
	}
}