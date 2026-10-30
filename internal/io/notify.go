package io

import (
	"fmt"
	"sync"

	"github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/common/util/jsonutil"
	"github.com/forfun/gforgame/internal/contract"
	"github.com/forfun/gforgame/internal/infra/net"
	"github.com/forfun/gforgame/internal/logger"
	"github.com/forfun/gforgame/internal/protos"
	serverconfig "github.com/forfun/gforgame/internal/serverconfig"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
	"github.com/forfun/gforgame/network/session"
)

var (
	notifyCodec          = json.NewSerializer()
	messageFactory       *protocol.MessageFactory
	gateSession          network.Session
	gateSessionMu        sync.RWMutex
	onlinePlayerRegistry *net.OnlinePlayerRegistry
)

// SetMessageFactory 设置业务协议注册表；在启动网络监听前调用。
func SetMessageFactory(factory *protocol.MessageFactory) {
	messageFactory = factory
}

func SetGateSession(session network.Session) {
	if session == nil {
		return
	}
	gateSessionMu.Lock()
	defer gateSessionMu.Unlock()
	gateSession = session
}

func SetOnlinePlayerRegistry(registry *net.OnlinePlayerRegistry) {
	onlinePlayerRegistry = registry
}

func NotifyByPlayerId(playerID string, index int32, data any) {
	if !serverconfig.GetCurrentNode().UseGateMode {
		// 直连模式：playerId 映射的是客户端会话，直接发送即可
		if s, found := session.GetSessionByUid(playerID); found {
			_ = s.Send(data, index)
		}
		return
	}

	// 网关模式：转发给 gate 服，再由 gate 回给客户端
	gs, ok := getGateSession()
	if !ok {
		logger.ErrorNoStack(fmt.Errorf("gate session not ready, player=%s", playerID))
		return
	}
	cmd, err := messageFactory.GetMessageCmd(data)
	if err != nil {
		return
	}
	body, err := notifyCodec.Encode(data)
	if err != nil {
		return
	}
	transferBody := &protos.TransferGateToLogic{
		PlayerId: playerID,
		Cmd:      cmd,
		Body:     body,
		Index:    index,
	}

	msgName, _ := messageFactory.GetMsgName(cmd)
	jsonStr, err := jsonutil.StructToJSON(data)
	logger.Info(fmt.Sprintf("id:%v 发送消息 cmd:%d, name:%s, 内容:%s", playerID, cmd, msgName, jsonStr))
	_ = gs.Send(transferBody, index)
}

func NotifyPlayer(player contract.Player, data any) {
	if player == nil || data == nil {
		return
	}
	playerID := player.GetId()
	NotifyByPlayerId(playerID, 0, data)
}

func getGateSession() (network.Session, bool) {
	gateSessionMu.RLock()
	defer gateSessionMu.RUnlock()
	if gateSession == nil {
		return nil, false
	}
	return gateSession, true
}
