// 混合服务，将一些小功能放到这里， 减少代码量
package mixture

import (
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"sync"
)

// 混合服务
type MixtureService struct{}

var (
	instance *MixtureService
	once     sync.Once
)

func GetMixtureService() *MixtureService {
	once.Do(func() {
		instance = &MixtureService{}
	})
	return instance
}

func (s *MixtureService) OnClientUploadEvent(player *player.Player, event int32)  {
	player.ExtendBox.ClientEvents[event]++
	context.EventBus.Publish(events.ClientDiyEvent, events.ClientCustomEvent{
		PlayerEvent: events.PlayerEvent{
			Player: player,
		},
		EventId: event,
	})
}