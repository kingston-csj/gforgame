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
		itemIdMap, err := util.ToIntIntMap(params, ";", "=")
		if err != nil {
			return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
		}
		for itemId, itemNum := range itemIdMap {
			item.GetItemService().AddByModelId(player, itemId, itemNum)
		}
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
		recharge.GetRechargeService().Recharge(player, rechargeId)
	case "add_scene_items":
		// add_scene_items 1001=1;1=2
		itemIdMap, err := util.ToIntIntMap(params, ";", "=")
		if err != nil {
			return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
		}
		for itemId, itemNum := range itemIdMap {
			err := item.GetSceneItemService().AddByModelId(player, itemId, itemNum)
			if err != nil {
				return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
			}
		}
	case "remove_scene_items":
		itemIdMap, err := util.ToIntIntMap(params, ";", "=")
		if err != nil {
			return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
		}
		for itemId, itemNum := range itemIdMap {
			item.GetSceneItemService().UseByModelId(player, itemId, itemNum)
		}
	}
	context.EventBus.Publish(events.PlayerEntityChange, player)

	return nil
}