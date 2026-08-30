package dispatch

import (
	"fmt"

	"github.com/forfun/gforgame/actor"
	"github.com/forfun/gforgame/common/logger"
)

// PlayerTaskDispatcher 用于在无玩家会话时，按 playerId 串行执行任务。
type PlayerTaskDispatcher struct {
	actorSystem *actor.ActorSystem
}

func NewPlayerTaskDispatcher(actorSystem *actor.ActorSystem) *PlayerTaskDispatcher {
	d := &PlayerTaskDispatcher{
		actorSystem: actorSystem,
	}
	return d
}

func (d *PlayerTaskDispatcher) DispatchPlayerTask(playerID string, task func()) {
	path := actor.NewActorPath("game", "player", playerID)
	// 查询actorRef
	ref := d.actorSystem.Find(path.String())
	if ref == nil {
		logger.ErrorNoStack(fmt.Errorf("player actor not found, playerId=%s", playerID))
		return
	}
	ref.Task(task)
}
