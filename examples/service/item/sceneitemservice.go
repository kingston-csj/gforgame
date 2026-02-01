package item

import (
	"io/github/gforgame/protos"

	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/service/catalog"

	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/reward"
	"io/github/gforgame/network"
)

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

func (s *SceneItemService) UseByModelId(p *playerdomain.Player, itemId int32, count int32) error {
	if itemId <= 0 || count <= 0 {
		return errorIllegalParams
	}
	backpack := p.SceneBackpack
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

func (s *SceneItemService) UseByUid(p *playerdomain.Player, itemUid string, count int32) (error, []protos.RewardVo) {
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
	return nil, make([]protos.RewardVo, 0)
}

func (s *SceneItemService) AddByModelId(p *playerdomain.Player, itemId int32, count int32) error {
	if itemId <= 0 || count <= 0 {
		return errorIllegalParams
	}

	itemData := config.QueryById[configdomain.ScenePropData](itemId)
	if itemData == nil {
		return errorIllegalParams
	}

	changeResult, err := p.SceneBackpack.AddByModelId(itemId, count)
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
