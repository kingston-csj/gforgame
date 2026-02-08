package mall

import (
	"io/github/gforgame/common"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/consume"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/reward"
	"sync"
)

// 商城模块
type MallService struct {
}

var (
	instance *MallService
	once     sync.Once
)

func GetMallService() *MallService {
	once.Do(func() {
		instance = &MallService{}
	})
	return instance
}

func (s *MallService) OnPlayerLogin(player *playerdomain.Player) {
}

func (s *MallService) Buy(player *playerdomain.Player, mallId int32, count int32) *common.BusinessRequestException {
	mallData := config.QueryById[configdomain.MallData](mallId)
	if mallData == nil {
		return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	andConsume := &consume.AndConsume{}
	for i := 0; i < int(count); i++ {
		andConsume.Add(consume.ParseConsume(mallData.Consumes))
	}
	andConsume = andConsume.Merge()
	err := andConsume.Verify(player)
	if err != nil {
		return err.(*common.BusinessRequestException)
	}
	// 每日限购
	if mallData.DailyBuy > 0 {
		// 检查是否超过每日限购
		if player.DailyReset.MallDailyBuy[mallId] >= mallData.DailyBuy {
			return common.NewBusinessRequestException(constants.I18N_MALL_DAILY_BUY_LIMIT)
		}
		// 增加每日购买次数
		player.DailyReset.MallDailyBuy[mallId]++
	}
	// 终身限购
	if mallData.LifeTimeBuy > 0 {
		// 检查是否超过终身限购
		if player.ExtendBox.LifeTimeBuyCount[mallId] >= mallData.LifeTimeBuy {
			return common.NewBusinessRequestException(constants.I18N_MALL_LIFE_TIME_BUY_LIMIT)
		}
		// 增加终身购买次数
		player.ExtendBox.LifeTimeBuyCount[mallId]++
	}
	// 消耗
	andConsume.Consume(player, constants.ActionType_BuyMall)
	// 奖励
	andReward := &reward.AndReward{}
	for i := 0; i < int(count); i++ {
		andReward.AddReward(reward.ParseReward(mallData.Rewards))
	}
	andReward = andReward.Merge()
	andReward.Reward(player, constants.ActionType_BuyMall)

	context.EventBus.Publish(events.MallBuy, &events.MallBuyEvent{
		PlayerEvent: events.PlayerEvent{
			Player: player,
		},
	})

	context.EventBus.Publish(events.PlayerEntityChange, player)
   
	return nil
}