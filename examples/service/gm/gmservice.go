package gm

import (
	"io/github/gforgame/common"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/reward"
	"io/github/gforgame/examples/service/item"
	questservice "io/github/gforgame/examples/service/quest"
	"io/github/gforgame/examples/service/recharge"
	"io/github/gforgame/logger"
	"io/github/gforgame/util"
	"strings"
	"sync"
)

type GmService struct{}

var (
	instance *GmService
	once     sync.Once
)

func GetGmService() *GmService {
	once.Do(func() {
		instance = &GmService{}
	})
	return instance
}

func (s *GmService) Dispatch(player *playerdomain.Player, topic string, params string) *common.BusinessRequestException{
	defer func() {
		if err := recover(); err != nil {
			logger.Error2("gm dispatch fail" , err.(error))
		}
	}()
	switch topic {
	case "add_items":
		itemParams := strings.Split(params, "=")
		itemId, err := util.StringToInt32(itemParams[0])
		if err != nil {
			return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
		}
		itemNum, err := util.StringToInt32(itemParams[1])
		if err != nil {
			return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
		}

		item.GetItemService().AddByModelId(player, itemId, itemNum)
	case "add_diamond":
		count, _ := util.StringToInt32(params)
		reward := &reward.CurrencyReward{
			Currency:   "diamond",
			Amount: count,
		}
		reward.Reward(player, constants.ActionType_Gm)
	case "add_gold":
		count, _ := util.StringToInt32(params)
		reward := &reward.CurrencyReward{
			Currency:   "gold",
			Amount: count,
		}
		reward.Reward(player, constants.ActionType_Gm)
	case "quest":
		questId, _ := util.StringToInt32(params)
		questservice.GetQuestService().GmFinish(player, questId)
	
	case "recharge":
		rechargeId, _ := util.StringToInt32(params)
		recharge.GetRechargeService().GmRecharge(player, rechargeId)
	}
	context.EventBus.Publish(events.PlayerEntityChange, player)

	return nil
}