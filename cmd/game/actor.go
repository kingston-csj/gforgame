package main

import (
	"github.com/forfun/gforgame/actor"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
)

type SharedAnonymousActor struct {
	router *network.MessageRoute
}

func NewSharedAnonymousActor(router *network.MessageRoute) *SharedAnonymousActor {
	return &SharedAnonymousActor{
		router: router,
	}
}

func (s *SharedAnonymousActor) OnStart() {
	logger.Info("shared anonymous actor started")
}

// PlayerMsgTask 投递到actor的消息任务，和前面定义一致
type PlayerMsgTask struct {
	session    network.Session
	frame      *protocol.RequestDataFrame
	msgHandler *network.Handler
}

// PlayerActor 玩家Actor，实现actor.Actor
type PlayerActor struct {
	playerId string
	router   *network.MessageRoute
}

func NewPlayerActor(playerId string, router *network.MessageRoute) *PlayerActor {
	return &PlayerActor{
		playerId: playerId,
		router:   router,
	}
}

func (p *PlayerActor) OnStart() {
	// 玩家Actor初始化，加载玩家数据、角色存档
}

func (p *PlayerActor) OnMessage(raw actor.Message) actor.Message {
	switch task := raw.(type) {
	case *PlayerMsgTask:
		// 在Actor单协程内，执行原来整套dispatch逻辑
		defer func() {
			if r := recover(); r != nil {
				handleRoutePanic(task.session, task.frame, task.msgHandler, "player actor handler panic", logger.PanicToError(r))
			}
		}()
		_ = dispatchMessage(task.session, task.frame, task.msgHandler)
	default:
	}
	return nil
}

func (p *PlayerActor) OnStop() {
	// 玩家下线，保存数据、清理资源
}

func (s *SharedAnonymousActor) OnMessage(raw actor.Message) actor.Message {
	switch task := raw.(type) {
	case *PlayerMsgTask:
		defer func() {
			if r := recover(); r != nil {
				handleRoutePanic(task.session, task.frame, task.msgHandler, "shared anon actor panic", logger.PanicToError(r))
			}
		}()
		// 复用原有消息分发逻辑，处理登录等匿名协议
		_ = dispatchMessage(task.session, task.frame, task.msgHandler)
	default:
	}
	return nil
}

func (s *SharedAnonymousActor) OnStop() {
	logger.Info("shared anonymous actor stopped")
}