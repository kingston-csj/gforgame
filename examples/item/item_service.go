package item

import (
	"io/github/gforgame/common"
	"io/github/gforgame/common/i18n"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/network"
	"sync"
)

type ItemService struct {
	network.Base
}

var instance *ItemService
var once sync.Once
var errorIllegalParams = common.NewBusinessRequestException(i18n.ErrorIllegalParams)

func GetInstance() *ItemService {
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

	itemData := context.GetDataManager().GetRecord("item", int64(itemId))
	if itemData == nil {
		return errorIllegalParams
	}

	p.Backpack.AddItem(itemId, count)

	return nil
}
