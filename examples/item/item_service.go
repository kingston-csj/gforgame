package item

import (
	"sync"

	"io/github/gforgame/common"
	"io/github/gforgame/protos"

	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/io"

	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/network"
)

type ItemService struct {
	network.Base
}

var (
	instance           *ItemService
	once               sync.Once
	errorIllegalParams = common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
)

var RecruitItemId int32 = 2002

func GetItemService() *ItemService {
	once.Do(func() {
		instance = &ItemService{}
		instance.init()
	})
	return instance
}

func (s *ItemService) init() {
}

func (s *ItemService) UseByModelId(p *playerdomain.Player, itemId int32, count int32) error {
	if itemId <= 0 || count <= 0 {
		return errorIllegalParams
	}

	return nil
}

func (s *ItemService) AddByModelId(p *playerdomain.Player, itemId int32, count int32) any {
	if itemId <= 0 || count <= 0 {
		return errorIllegalParams
	}

	itemData := config.QueryById[configdomain.ItemData](int64(itemId))
	if itemData == nil {
		return errorIllegalParams
	}

	p.Backpack.AddItem(itemId, count)

	context.EventBus.Publish(events.PlayerEntityChange, p)

	io.NotifyPlayer(p, &protos.PushItemChanged{
		ItemId: itemId,
		Count:  count,
	})

	return nil
}
