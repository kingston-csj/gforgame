package cross

import (
	"fmt"

	"github.com/forfun/gforgame/common/logger"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
)

type IslandTransfer struct{}

func init() {
	RegisterTransfer(Island, &IslandTransfer{})
}

func (t IslandTransfer) CanTransfer(p *playerdomain.Player) int {
	return 0
}

func (t IslandTransfer) GetTargetServerId(p *playerdomain.Player) uint32 {
	return 0
}

func (t IslandTransfer) LocalEnterScene(p *playerdomain.Player) error {
	logger.Info(fmt.Sprintf("player %s enter island", p.Id))
	return nil
}

func (t IslandTransfer) RemoteEnterScene(p *playerdomain.Player) error {
	return nil
}
