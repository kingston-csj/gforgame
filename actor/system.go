package actor

import (
	"fmt"
	"strings"
	"sync"

	"github.com/forfun/gforgame/common/logger"
)

// ActorSystem actor容器系统
type ActorSystem struct {
	mu     sync.RWMutex
	actors map[string]*ActorRef // key = path.String()
	mailboxCap int // 默认256
}

func NewActorSystem() *ActorSystem {
	return &ActorSystem{
		actors: make(map[string]*ActorRef),
		mailboxCap: 256,
	}
}

// Spawn 创建并启动一个actor
func (sys *ActorSystem) Spawn(path ActorPath, actor Actor) (*ActorRef, error) {
	pathStr := path.String()

	sys.mu.Lock()
	defer sys.mu.Unlock()
	if ref, ok := sys.actors[pathStr]; ok {
		// 已存在，直接返回
		return ref, nil
	}

	mb := NewMailbox(sys.mailboxCap)
	ref := &ActorRef{
		path:    path,
		mailbox: mb,
	}
	sys.actors[pathStr] = ref

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorNoStack(fmt.Errorf("actor %s panic recovered: %v", pathStr, r))
			}
		}()
		actor.OnStart()
		for raw := range mb.Recv() {
			switch v := raw.(type) {
			case askWrap:
				ret := actor.OnMessage(v.payload)
				// 关键bug修复：防止对方放弃接收resp，卡死整个actor协程
				select {
				case v.resp <- ret:
				default:
					// resp通道已经没人接收，直接丢弃返回值，不阻塞主循环
				}
			case taskMsg:
				// 直接执行task闭包，不需要调用OnMessage
				v.fn()	
			default:
				_ = actor.OnMessage(raw)
			}
		}
		// 通道关闭，队列全部消息处理完毕，执行停止回调
		actor.OnStop()
		// 退出时，从系统注册表删除自己
		sys.mu.Lock()
		delete(sys.actors, pathStr)
		sys.mu.Unlock()
	}()

	return ref, nil
}


// GetOrCreate 原子获取或创建Actor
// builder: 当actor不存在时才会调用，只调用一次
func (sys *ActorSystem) GetOrCreate(path ActorPath, builder func() (Actor, error)) (*ActorRef, error) {
	pathStr := path.String()
	sys.mu.Lock()
	defer sys.mu.Unlock()

	// 已经存在直接返回
	if ref, ok := sys.actors[pathStr]; ok {
		return ref, nil
	}

	// 不存在，调用builder构造actor实例
	actor, err := builder()
	if err != nil {
		return nil, err
	}

	mb := NewMailbox(sys.mailboxCap)
	ref := &ActorRef{
		path:    path,
		mailbox: mb,
	}
	sys.actors[pathStr] = ref

	go func() {
		actor.OnStart()
		for raw := range mb.Recv() {
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
		actor.OnStop()
		// 退出删除注册表
		sys.mu.Lock()
		delete(sys.actors, pathStr)
		sys.mu.Unlock()
	}()

	return ref, nil
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
