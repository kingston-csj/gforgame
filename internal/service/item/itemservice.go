package item

import (
	"github.com/forfun/gforgame/common/errors"
	"github.com/forfun/gforgame/internal/protos"

	"github.com/forfun/gforgame/internal/config"
	"github.com/forfun/gforgame/internal/constants"
	"github.com/forfun/gforgame/internal/consume"
	"github.com/forfun/gforgame/internal/context"
	"github.com/forfun/gforgame/internal/contract"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/service/catalog"
	playerservice "github.com/forfun/gforgame/internal/service/player"

	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/reward"
	"github.com/forfun/gforgame/network"
)

// 普通道具模块
type ItemService struct {
	network.Base
	player  *playerservice.PlayerService
	catalog *catalog.CatalogService
}

var (
	itemservice           *ItemService
	errorIllegalParams = errors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
	notEnoughError = errors.NewBusinessError(constants.I18N_ITEM_NOT_ENOUGH)
)

var RecruitItemId int32 = 2002

func NewItemService(player *playerservice.PlayerService, catalogService *catalog.CatalogService) *ItemService {
	service := &ItemService{
		player:  player,
		catalog: catalogService,
	}
	service.init()
	return service
}

func (s *ItemService) init() {
	reward.SetItemOps( s)
	consume.SetItemOps(s)
}

func (s *ItemService) UseByModelId(playerId string, itemId int32, count int32) error {
	p := s.player.GetPlayer(playerId)
	backpack := p.Backpack
	if itemId <= 0 || count <= 0 {
		return errorIllegalParams
	}
	changeResult := backpack.ReduceByModelId(itemId, count)
	if !changeResult.Succ  {
		return notEnoughError
	}

	context.EventBus.Publish(events.ItemConsume, &events.ItemConsumeEvent{
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

func (s *ItemService) UseByUid(p *playerdomain.Player, itemUid string, count int32) (error, []contract.RewardDefLite) {
	if itemUid == "" || count <= 0 {
		return errorIllegalParams, nil
	}

	return nil, nil
}

func (s *ItemService) AddByModelId(playerId string, itemId int32, count int32) error {
	p := s.player.GetPlayer(playerId)
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
	s.catalog.TryUnlock(p, 1, itemId)

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
