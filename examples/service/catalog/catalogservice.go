package catalog

import (
	"sync"

	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/reward"
	"io/github/gforgame/protos"
)

type CatalogService struct {
}

var (
	instance *CatalogService
	once     sync.Once
)

func GetCatalogService() *CatalogService {
	once.Do(func() {
		instance = &CatalogService{}
	})
	return instance
}

func (s *CatalogService) OnPlayerLogin(player *playerdomain.Player) {
	push := &protos.PushCatalogInfo{
		SitemCatalog: buildCatalogVo(&player.ExtendBox.SitemCatalogModel),
		ItemCatalog:  buildCatalogVo(&player.ExtendBox.ItemCatalogModel),
		MenuCatalog:  buildCatalogVo(&player.ExtendBox.MenuCatalogModel),
	}
	io.NotifyPlayer(player, push)
}

func buildCatalogVo(catalogModel *playerdomain.CatalogModel) protos.CatalogModel {
	return protos.CatalogModel{
		UnlockIds:   catalogModel.UnlockIds.ToSlice(),
		ReceivedIds: catalogModel.ReceivedIds.ToSlice(),
	}
}

func (s *CatalogService) TakeReward(player *playerdomain.Player, typ int32, id int32) (int, []*protos.RewardVo) {
	catalogModel := getcatalogModel(player, typ)
	rewardStr := ""
	if typ == 0 {
		itemData := config.QueryById[configdomain.ScenePropData](id)
		rewardStr = itemData.ActivateRewards
	} else if typ == 1 {
		itemData := config.QueryById[configdomain.PropData](id)
		rewardStr = itemData.ActivateRewards
	} else {
		itemData := config.QueryById[configdomain.MenuData](id)
		rewardStr = itemData.ActivateRewards
	}
	if !catalogModel.CanReceived(id) {
		 return constants.I18N_COMMON_ILLEGAL_PARAMS,nil
	}
	rewards := reward.ParseReward(rewardStr)
	rewards.Reward(player, constants.ActionType_CatalogActivate)
	return 0, reward.ToRewardVos(rewards)
}

func (s *CatalogService) TryUnlock(player *playerdomain.Player, typ int32, id int32) bool {
	catalogModel := getcatalogModel(player, typ)
	succ := catalogModel.AddUnlock(id)
	if succ {
		push := &protos.PushCatalogAdd{
			Typ: typ,
			ItemId:  id,
		}
		io.NotifyPlayer(player, push)
	}
	return succ
}

func getcatalogModel(player *playerdomain.Player, typ int32) *playerdomain.CatalogModel {
	if typ == 0 {
		return &player.ExtendBox.SitemCatalogModel
	} else if typ == 1 {
		return &player.ExtendBox.ItemCatalogModel
	} else {
		return &player.ExtendBox.MenuCatalogModel
	}
}