package io

import (
	"fmt"
	"sync"

	"github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/common/logger"
	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/internal/contract"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/network"
)

var (
	notifyCodec   = json.NewSerializer()
	gateSession   *network.Session
	gateSessionMu sync.RWMutex
)

func SetGateSession(session *network.Session) {
	if session == nil {
		return
	}
	gateSessionMu.Lock()
	defer gateSessionMu.Unlock()
	gateSession = session
}

func NotifyPlayer(player contract.Player, data any) {
	if player == nil || data == nil {
		return
	}
	playerID := player.GetId()

	if !serverconfig.ServerConfig.UseGateMode {
		// 直连模式：playerId 映射的是客户端会话，直接发送即可
		if s := network.GetSessionByPlayerId(playerID); s != nil {
			_ = s.SendWithoutIndex(data)
		}
		return
	}

	// 网关模式：转发给 gate 服，再由 gate 回给客户端
	gs, ok := getGateSession()
	if !ok {
		logger.ErrorNoStack(fmt.Errorf("gate session not ready, player=%s", playerID))
		return
	}
	cmd, err := network.GetMessageCmd(data)
	if err != nil {
		return
	}
	body, err := notifyCodec.Encode(data)
	if err != nil {
		return
	}
	if err := gs.SendWithoutIndex(&protos.TransferGateToLogic{
		PlayerId: playerID,
		Cmd:      cmd,
		Body:     body,
	}); err != nil {
		return
	}
}

func getGateSession() (*network.Session, bool) {
	gateSessionMu.RLock()
	defer gateSessionMu.RUnlock()
	if gateSession == nil {
		return nil, false
	}
	return gateSession, true
}
