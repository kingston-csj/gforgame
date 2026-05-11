package recharge

import (
	"strings"
	"sync"
	"time"

	"github.com/forfun/gforgame/common/util/conv"
	config "github.com/forfun/gforgame/internal/config"
	"github.com/forfun/gforgame/internal/constants"
	"github.com/forfun/gforgame/internal/context"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/reward"
)

// 充值模块
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

func (s *RechargeService) Recharge(player *playerdomain.Player, rechargeId int32) {
	rechargeData := config.QueryById[configdomain.RechargeData](rechargeId)	
	if conv.IsEmptyString(rechargeData.Children) {
		s.recharge0(player, rechargeId)
	} else {
		// 如果是捆绑销售的商品
		children := strings .Split(rechargeData.Children, ",")
		for _, child := range children {
			s.recharge0(player, conv.Int32Value(child))
		}
	}
}

func (s *RechargeService) recharge0(player *playerdomain.Player, rechargeId int32) {
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
	evt := &events.RechargeEvent{
		PlayerEvent: events.PlayerEvent{
			Player: player,
		},
		RechargeId: rechargeId,
	}
	context.EventBus.Publish(events.Recharge, evt)
}