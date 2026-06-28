package mall

import (
	"github.com/forfun/gforgame/common/errors"
	"github.com/forfun/gforgame/internal/config"
	"github.com/forfun/gforgame/internal/constants"
	"github.com/forfun/gforgame/internal/consume"
	"github.com/forfun/gforgame/internal/context"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/reward"
)

// 商城模块
type MallService struct {
}

func NewMallService() *MallService {
	return &MallService{}
}

func (s *MallService) OnPlayerLogin(player *playerdomain.Player) {
}

func (s *MallService) Buy(player *playerdomain.Player, mallId int32, count int32) *errors.BusinessError {
	mallData := config.QueryById[configdomain.MallData](mallId)
	if mallData == nil {
		return errors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	andConsume := &consume.AndConsume{}
	for i := 0; i < int(count); i++ {
		andConsume.Add(consume.ParseConsume(mallData.Consumes))
	}
	andConsume = andConsume.Merge()
	err := andConsume.Verify(player)
	if err != nil {
		return err.(*errors.BusinessError)
	}
	// 每日限购
	if mallData.DailyBuy > 0 {
		// 检查是否超过每日限购
		if player.DailyReset.MallDailyBuy[mallId] >= mallData.DailyBuy {
			return errors.NewBusinessError(constants.I18N_MALL_DAILY_BUY_LIMIT)
		}
		// 增加每日购买次数
		player.DailyReset.MallDailyBuy[mallId]++
	}
	// 终身限购
	if mallData.LifeTimeBuy > 0 {
		// 检查是否超过终身限购
		if player.ExtendBox.LifeTimeBuyCount[mallId] >= mallData.LifeTimeBuy {
			return errors.NewBusinessError(constants.I18N_MALL_LIFE_TIME_BUY_LIMIT)
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