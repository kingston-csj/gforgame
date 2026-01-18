package mall

import (
	"io/github/gforgame/common"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"sync"
)

type MallService struct {
}

var (
	instance *MallService
	once     sync.Once
)

func GetMallService() *MallService {
	once.Do(func() {
		instance = &MallService{}
	})
	return instance
}

func (s *MallService) OnPlayerLogin(player *playerdomain.Player) {
}

func (s *MallService)  Buy(player *playerdomain.Player, mallId int32, count int32) *common.BusinessRequestException {

		mallData := config.QueryById[configdomain.MallData](mallId)
	if mallData == nil {
		return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
   
	return nil
}