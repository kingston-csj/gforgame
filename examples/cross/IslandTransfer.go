package cross

import (
	"fmt"
	"io/github/gforgame/examples/player"
	"io/github/gforgame/log"
)

type IslandTransfer struct {
}

func init() {
	RegisterTransfer(Island, &IslandTransfer{})
}

func (t IslandTransfer) CanTransfer(p *player.Player) int {
	return 0
}

func (t IslandTransfer) GetTargetServerId(p *player.Player) uint32 {
	return 0
}

func (t IslandTransfer) LocalEnterScene(p *player.Player) error {
	log.Info(fmt.Sprintf("player %s enter island", p.Id))
	return nil
}

func (t IslandTransfer) RemoteEnterScene(p *player.Player) error {
	return nil
}
