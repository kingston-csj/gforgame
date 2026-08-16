package dispatch

import (
	"fmt"
	"hash/fnv"

	"github.com/forfun/gforgame/common/logger"
	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/internal/infra/net"
)

// PlayerTaskDispatcher 用于在无玩家会话时，按 playerId 串行执行任务。
type PlayerTaskDispatcher struct {
	workerCount uint32
	queues      []chan func()

	onlinePlayerRegistry *net.OnlinePlayerRegistry
}

func NewPlayerTaskDispatcher(workerCount uint32, onlinePlayerRegistry *net.OnlinePlayerRegistry) *PlayerTaskDispatcher {
	if workerCount == 0 {
		workerCount = 1
	}
	d := &PlayerTaskDispatcher{
		workerCount:          workerCount,
		queues:               make([]chan func(), workerCount),
		onlinePlayerRegistry: onlinePlayerRegistry,
	}
	for i := uint32(0); i < workerCount; i++ {
		q := make(chan func(), 512)
		d.queues[i] = q
		go func(ch chan func()) {
			for task := range ch {
				func() {
					defer func() {
						if r := recover(); r != nil {
							logger.Error("player task dispatcher: panic recovered: %v", fmt.Errorf("player task panic: %v", r))
						}
					}()
					task()
				}()
			}
		}(q)
	}
	return d
}

func (d *PlayerTaskDispatcher) submit(playerID string, task func()) {
	idx := hashPlayerID(playerID) % d.workerCount
	d.queues[idx] <- task
}

func hashPlayerID(playerID string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(playerID))
	return h.Sum32()
}

// DispatchPlayerTask 保证同一玩家任务线程安全
// 1. 网关模式：统一走全局 playerId 分片队列，避免共享会话导致全玩家串行
// 2. 直连模式：优先投递到玩家会话 AsynTasks
func (d *PlayerTaskDispatcher) DispatchPlayerTask(playerID string, task func()) {
	if playerID == "" || task == nil {
		return
	}
	if serverconfig.ServerConfig.UseGateMode {
		d.submit(playerID, task)
		return
	}
	if session := d.onlinePlayerRegistry.GetSessionByPlayerID(playerID); session != nil {
		select {
		case <-session.DieChan():
		default:
			session.AsynTasksChan() <- task
			return
		}
	}
	d.submit(playerID, task)
}
