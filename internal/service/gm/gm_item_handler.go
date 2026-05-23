package gm

import (
	commonerrors "github.com/forfun/gforgame/common/errors"
	"github.com/forfun/gforgame/common/util/conv"
	"github.com/forfun/gforgame/internal/constants"
	"github.com/forfun/gforgame/internal/consume"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/reward"
	"github.com/forfun/gforgame/internal/service/item"
)

type ItemGmHandler struct {
	item      *item.ItemService
}

func NewItemGmHandler(itemService *item.ItemService) *ItemGmHandler {
	return &ItemGmHandler{
		item:      itemService,

	}
}

func (h *ItemGmHandler) RegisterTo(gm *GmService) {
	gm.Register("add_items", "添加物品", "add_items 1001=1;1002=2", h.handleAddItems)
	gm.Register("remove_items", "移除物品", "remove_items 1001=1", h.handleRemoveItems)
	gm.Register("add_diamond", "添加钻石", "add_diamond 1000", h.handleAddDiamond)
	gm.Register("add_gold", "添加金币", "add_gold 1000", h.handleAddGold)
}

func (h *ItemGmHandler) handleAddItems(player *playerdomain.Player, params string) *commonerrors.BusinessError {
	itemIdMap, err := conv.ToIntIntMap(params, ";", "=")
	if err != nil {
		return commonerrors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	for itemId, itemNum := range itemIdMap {
		h.item.AddByModelId(player.Id, itemId, itemNum)
	}
	return nil
}


func (h *ItemGmHandler) handleRemoveItems(player *playerdomain.Player, params string) *commonerrors.BusinessError {
	itemIdMap, err := conv.ToIntIntMap(params, ";", "=")
	if err != nil {
		return commonerrors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	consums := &consume.AndConsume{}
	for itemId, itemNum := range itemIdMap {
		consums.Add(&consume.ItemConsume{
			ItemId: itemId,
			Amount: itemNum,
		})
	}
	if err := consums.Verify(player); err != nil {
		return err.(*commonerrors.BusinessError)
	}
	consums.Consume(player, constants.ActionType_Gm)
	return nil
}

func (h *ItemGmHandler) handleAddDiamond(player *playerdomain.Player, params string) *commonerrors.BusinessError {
	count, _ := conv.StringToInt32(params)
	reward := &reward.CurrencyReward{
		Currency: "diamond",
		Amount:   count,
	}
	reward.Reward(player, constants.ActionType_Gm)
	return nil
}

func (h *ItemGmHandler) handleAddGold(player *playerdomain.Player, params string) *commonerrors.BusinessError {
	count, _ := conv.StringToInt32(params)
	reward := &reward.CurrencyReward{
		Currency: "gold",
		Amount:   count,
	}
	reward.Reward(player, constants.ActionType_Gm)
	return nil
}

