package network

import (
	"hash/fnv"
	"log/slog"

	serverconfig "github.com/forfun/gforgame/config"
)

// playerTaskDispatcher 用于在无玩家会话时，按 playerId 串行执行任务。
type playerTaskDispatcher struct {
	workerCount uint32
	queues      []chan func()
}

var globalPlayerTaskDispatcher = newPlayerTaskDispatcher(32)

func newPlayerTaskDispatcher(workerCount uint32) *playerTaskDispatcher {
	if workerCount == 0 {
		workerCount = 1
	}
	d := &playerTaskDispatcher{
		workerCount: workerCount,
		queues:      make([]chan func(), workerCount),
	}
	for i := uint32(0); i < workerCount; i++ {
		q := make(chan func(), 512)
		d.queues[i] = q
		go func(ch chan func()) {
			for task := range ch {
				func() {
					defer func() {
						if r := recover(); r != nil {
							slog.Error("player task panic", "panic", r)
						}
					}()
					task()
				}()
			}
		}(q)
	}
	return d
}

func (d *playerTaskDispatcher) submit(playerID string, task func()) {
	idx := hashPlayerID(playerID) % d.workerCount
	d.queues[idx] <- task
}

func hashPlayerID(playerID string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(playerID))
	return h.Sum32()
}

// DispatchPlayerTask 保证同一玩家任务串行：
// 1. 网关模式：统一走全局 playerId 分片队列，避免共享会话导致全玩家串行
// 2. 直连模式：优先投递到玩家会话 AsynTasks（保持原有单玩家执行语义）
func DispatchPlayerTask(playerID string, task func()) {
	if playerID == "" || task == nil {
		return
	}
	if serverconfig.ServerConfig.UseGateMode {
		globalPlayerTaskDispatcher.submit(playerID, task)
		return
	}
	if session := GetSessionByPlayerId(playerID); session != nil {
		select {
		case <-session.Die:
		default:
			session.AsynTasks <- task
			return
		}
	}
	globalPlayerTaskDispatcher.submit(playerID, task)
}
