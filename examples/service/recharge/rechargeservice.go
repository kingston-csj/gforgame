package recharge

import (
	config "io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/reward"
	"io/github/gforgame/protos"
	"sync"
	"time"
)


type RechargeService struct {
}

var (
	instance *RechargeService
	once     sync.Once
)

func GetRechargeService() *RechargeService {
	return &RechargeService{}
}


func (s *RechargeService) OnPlayerLogin(player *playerdomain.Player) {
    push := &protos.PushRechargeInfo{};
	io.NotifyPlayer(player, push)
}

func (s *RechargeService) GmRecharge(player *playerdomain.Player, rechargeId int32) {
	rechargeData := config.QueryById[configdomain.RechargeData](rechargeId)
	rewards := reward.ParseReward(rechargeData.Rewards)
	rewards.Reward(player, constants.ActionType_Recharge)

	push := &protos.PushRechargePay{
		RechargeId: rechargeId,
		Rewards: reward.ToRewardVos(rewards),
	}
	io.NotifyPlayer(player, push)

	historyTimes := player.RechargeBox.RechargeTimes[rechargeId] 
	if historyTimes == 0 {
		player.RechargeBox.RechargeTimes[rechargeId] = 1
	} else {
		player.RechargeBox.RechargeTimes[rechargeId] = historyTimes + 1
	}

	player.DailyReset.RechargeSum += rechargeData.Money
	player.WeeklyReset.RechargeSum += rechargeData.Money
	player.MonthlyReset.RechargeSum += rechargeData.Money

	switch rechargeData.Type {
	case constants.RechargeTypeFirst:
		 player.RechargeBox.FirstRechargeTime = time.Now().Unix()
	}

	context.EventBus.Publish(events.PlayerEntityChange, player)
	context.EventBus.Publish(events.Recharge, rechargeId)
}