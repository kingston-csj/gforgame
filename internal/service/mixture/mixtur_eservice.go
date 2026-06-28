// 混合服务，将一些小功能放到这里， 减少代码量
package mixture

import (
	"github.com/forfun/gforgame/internal/context"
	"github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
)

// 混合服务
type MixtureService struct{}

func NewMixtureService() *MixtureService {
	return &MixtureService{}
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