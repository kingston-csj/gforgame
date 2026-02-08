package item

import (
	"sync"

	"io/github/gforgame/common"
	"io/github/gforgame/protos"

	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/consume"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/service/catalog"

	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/reward"
	"io/github/gforgame/network"
)

// 普通道具模块
type ItemService struct {
	network.Base
}

var (
	itemservice           *ItemService
	once               sync.Once
	errorIllegalParams = common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	notEnoughError = common.NewBusinessRequestException(constants.I18N_ITEM_NOT_ENOUGH)
)

var RecruitItemId int32 = 2002

func GetItemService() *ItemService {
	once.Do(func() {
		itemservice = &ItemService{}
		itemservice.init()
	})
	return itemservice
}

func (s *ItemService) init() {
	reward.SetItemOps( s)
	consume.SetItemOps(s)
}

func (s *ItemService) UseByModelId(p *playerdomain.Player, itemId int32, count int32) error {
	if itemId <= 0 || count <= 0 {
		return errorIllegalParams
	}
	backpack := p.Backpack
	changeResult := backpack.ReduceByModelId(itemId, count)
	if !changeResult.Succ  {
		return notEnoughError
	}

	context.EventBus.Publish(events.ItemConsume, events.ItemConsumeEvent{
		PlayerEvent: events.PlayerEvent{
			Player: p,
		},
		ItemId: itemId,
		Count:  count,
	})

	notify :=  &protos.PushItemChanged{
		Type: "item",
		Changed: changeResult.ToChangeInfos(),
	}
	io.NotifyPlayer(p, notify)

	return nil
}

func (s *ItemService) UseByUid(p *playerdomain.Player, itemUid string, count int32) (error, []protos.RewardVo) {
	if itemUid == "" || count <= 0 {
		return errorIllegalParams, nil
	}

	return nil, nil
}

func (s *ItemService) AddByModelId(p *playerdomain.Player, itemId int32, count int32) error {
	if itemId <= 0 || count <= 0 {
		return errorIllegalParams
	}

	itemData := config.QueryById[configdomain.PropData](itemId)
	if itemData == nil {
		return errorIllegalParams
	}

	changeResult, err := p.Backpack.AddByModelId(itemId, count, func(item *playerdomain.Item) {
		item.Type = constants.BackpackType_Norm
	})
	if err != nil {
		return err
	}
	// 激活图鉴
	catalog.GetCatalogService().TryUnlock(p, 1, itemId)

	// 发布事件，供任务系统使用
	context.EventBus.Publish(events.PlayerEntityChange, p)
	
	itemInfos := make([]protos.ItemInfo, 0, len(changeResult.ChangedItems))
	for _, item := range changeResult.ChangedItems {
		itemInfos = append(itemInfos, item.Item.ToVo())
	}

	notify := &protos.PushItemChanged{
		Type: "item",
		Changed: itemInfos,
	}
	io.NotifyPlayer(p, notify)
	
	return nil
}
