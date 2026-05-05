package item

import (
	"github.com/forfun/gforgame/examples/protos"

	"github.com/forfun/gforgame/examples/config"
	"github.com/forfun/gforgame/examples/constants"
	"github.com/forfun/gforgame/examples/context"
	"github.com/forfun/gforgame/examples/contract"
	configdomain "github.com/forfun/gforgame/examples/domain/config"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	"github.com/forfun/gforgame/examples/events"
	"github.com/forfun/gforgame/examples/io"
	"github.com/forfun/gforgame/examples/reward"
	"github.com/forfun/gforgame/examples/service/catalog"
	playerservice "github.com/forfun/gforgame/examples/service/player"
	"github.com/forfun/gforgame/network"
)

// 场景道具模块
type SceneItemService struct {
	network.Base
}

var (
	sceneItemService           *SceneItemService
)

func GetSceneItemService() *SceneItemService {
	once.Do(func() {
		sceneItemService = &SceneItemService{}
		sceneItemService.init()
	})
	return sceneItemService
}

func (s *SceneItemService) init() {
	reward.SetSceneItemOps( s)
}

func (s *SceneItemService) UseByModelId(playerId string, itemId int32, count int32) error {
	p := playerservice.GetPlayerService().GetPlayer(playerId)
	backpack := p.SceneBackpack
	if itemId <= 0 || count <= 0 {
		return errorIllegalParams
	}
	changeResult := backpack.ReduceByModelId(itemId, count)
	if !changeResult.Succ  {
		return notEnoughError
	}

	notify :=  &protos.PushItemChanged{
		Type: "sitem",
		Changed: changeResult.ToChangeInfos(),
	}
	io.NotifyPlayer(p, notify)

	return nil
}

func (s *SceneItemService) UseByUid(playerId string, itemUid string, count int32) (error, []contract.RewardDefLite) {
	p := playerservice.GetPlayerService().GetPlayer(playerId)
	if itemUid == "" || count <= 0 {
		return errorIllegalParams, nil
	}
	backpack := p.SceneBackpack
	changeResult, err := backpack.ReduceByUid(itemUid, count)
	if err != nil {
		return err, nil
	}
	if !changeResult.Succ  {
		return notEnoughError, nil
	}

	notify :=  &protos.PushItemChanged{
		Type: "sitem",
		Changed: changeResult.ToChangeInfos(),
	}
	io.NotifyPlayer(p, notify)
	// 场景道具没有使用效果，直接返回空奖励
	return nil, make([]contract.RewardDefLite, 0)
}

func (s *SceneItemService) AddByModelId(playerId string, itemId int32, count int32) error {
	p := playerservice.GetPlayerService().GetPlayer(playerId)
	if itemId <= 0 || count <= 0 {
		return errorIllegalParams
	}

	itemData := config.QueryById[configdomain.ScenePropData](itemId)
	if itemData == nil {
		return errorIllegalParams
	}

	changeResult, err := p.SceneBackpack.AddByModelId(itemId, count, func(item *playerdomain.Item) {
		item.Type =  constants.BackpackType_SItem
		item.Extra = "0"
	})
	if err != nil {
		return err
	}

	// 激活图鉴
	catalog.GetCatalogService().TryUnlock(p, 0, itemId)
	context.EventBus.Publish(events.PlayerEntityChange, p)
	
	itemInfos := make([]protos.ItemInfo, 0, len(changeResult.ChangedItems))
	for _, item := range changeResult.ChangedItems {
		itemInfos = append(itemInfos, item.Item.ToVo())
	}

	notify := &protos.PushItemChanged{
		Type: "sitem",
		Changed: itemInfos,
	}
	io.NotifyPlayer(p, notify)
	
	return nil
}
