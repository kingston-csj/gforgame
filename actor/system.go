package actor

import (
	"fmt"
	"strings"
	"sync"

	"github.com/forfun/gforgame/common/logger"
)

// ActorSystem actor容器系统
type ActorSystem struct {
	mu         sync.RWMutex
	actors     map[string]*ActorRef // key = path.String()
	mailboxCap int                  // 默认256
}

// newActorRef 统一创建 ActorRef，避免两套入口各自组装邮箱。
func (sys *ActorSystem) newActorRef(path ActorPath) *ActorRef {
	return &ActorRef{
		path:    path,
		mailbox: NewMailbox(sys.mailboxCap),
	}
}

// runActorLoop 统一启动 actor 主循环，复用单条消息隔离和退出清理逻辑。
func (sys *ActorSystem) runActorLoop(pathStr string, ref *ActorRef, actor Actor, panicHandler func(any)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicHandler(r)
			}
		}()
		actor.OnStart()
		for raw := range ref.mailbox.Recv() {
			handleMessageSafely(pathStr, ref, actor, raw)
		}
		actor.OnStop()
		sys.mu.Lock()
		delete(sys.actors, pathStr)
		sys.mu.Unlock()
	}()
}

// handleMessageSafely 兜住单条消息的 panic，避免打死整个 actor 协程。
func handleMessageSafely(pathStr string, ref *ActorRef, actor Actor, raw Message) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("actor message panic", fmt.Errorf("actor message panic path=%s ref=%p msgType=%T err=%v", pathStr, ref, raw, r))
		}
	}()

	switch v := raw.(type) {
	case askWrap:
		ret := actor.OnMessage(v.payload)
		select {
		case v.resp <- ret:
		default:
		}
	case taskMsg:
		v.fn()
	default:
		_ = actor.OnMessage(raw)
	}
}

func NewActorSystem() *ActorSystem {
	return &ActorSystem{
		actors:     make(map[string]*ActorRef),
		mailboxCap: 256,
	}
}

// Spawn 创建并启动一个actor
func (sys *ActorSystem) Spawn(path ActorPath, actor Actor) *ActorRef {
	pathStr := path.String()

	sys.mu.Lock()
	defer sys.mu.Unlock()
	if ref, ok := sys.actors[pathStr]; ok {
		// 已存在，直接返回
		return ref
	}

	ref := sys.newActorRef(path)
	sys.actors[pathStr] = ref
	sys.runActorLoop(pathStr, ref, actor, func(r any) {
		logger.ErrorNoStack(fmt.Errorf("actor %s panic recovered: %v", pathStr, r))
	})

	return ref
}

// GetOrCreate 原子获取或创建Actor
// builder: 当actor不存在时才会调用，只调用一次
func (sys *ActorSystem) GetOrCreate(path ActorPath, builder func() (Actor)) *ActorRef {
	pathStr := path.String()
	sys.mu.Lock()
	defer sys.mu.Unlock()

	// 已经存在直接返回
	if ref, ok := sys.actors[pathStr]; ok {
		return ref
	}

	// 不存在，调用builder构造actor实例
	actor := builder()

	ref := sys.newActorRef(path)
	sys.actors[pathStr] = ref
	sys.runActorLoop(pathStr, ref, actor, func(r any) {
		logger.Error("actor run error", fmt.Errorf("actor run error path=%s ref=%p err=%v", pathStr, ref, r))
	})

	return ref
}

// Stop 外部停止actor
func (sys *ActorSystem) Stop(pathStr string) {
	sys.mu.RLock()
	ref, ok := sys.actors[pathStr]
	sys.mu.RUnlock()
	if ok {
		ref.Stop()
	}
}

// Find 根据完整路径字符串查找ActorRef
func (sys *ActorSystem) Find(pathStr string) *ActorRef {
	sys.mu.RLock()
	defer sys.mu.RUnlock()
	return sys.actors[pathStr]
}

// FindByPrefix 前缀过滤查询，例如 "game/player/"
func (sys *ActorSystem) FindByPrefix(prefix string) []*ActorRef {
	sys.mu.RLock()
	defer sys.mu.RUnlock()
	var list []*ActorRef
	for k, v := range sys.actors {
		if strings.HasPrefix(k, prefix) {
			list = append(list, v)
		}
	}
	return list
}
