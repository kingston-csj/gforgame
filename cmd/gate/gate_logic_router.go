package main

import (
	"fmt"

	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/conv"
	"github.com/forfun/gforgame/gateway/contract"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
)

// 作为逻辑层，接收logic层的推送
type LogicRouter struct {
	router *network.MessageRoute
}

func (g *LogicRouter) MessageReceived(session *network.Session, frame *protocol.RequestDataFrame) bool {
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorNoStack(fmt.Errorf("panic recovered: %v", r))
		}
	}()
	if frame.Header.Cmd == gateTransferCodec.TransferCmd() {
		transferResp, err := gateTransferCodec.ParseTransferMessage(frame.Msg)
		if err != nil {
			logger.ErrorNoStack(fmt.Errorf("decode transfer response failed: %v", err))
			return false
		}
		if err := forwardTransferToClient(session, transferResp); err != nil {
			logger.ErrorNoStack(err)
			return false
		}
		return true
	}
	return true
}

func forwardTransferToClient(logicSession *network.Session, transfer contract.GateTransferMessage) error {
	playerID := transfer.GetPlayerID()
	cmd := transfer.GetTransferCmd()
	index := transfer.GetTransferIndex()
	body := transfer.GetTransferBody()
	if playerID == "" {
		return fmt.Errorf("transfer response playerId is empty, cmd=%d", cmd)
	}
	serverID := resolveBackendServerID(logicSession)
	if serverID <= 0 {
		return fmt.Errorf("logic session serverId is empty, playerId=%s cmd=%d", playerID, cmd)
	}
	sessionPlayerKey := buildSessionPlayerKey(serverID, playerID)
	clientSession := network.GetSessionByPlayerId(sessionPlayerKey)
	if clientSession == nil {
		return fmt.Errorf("client session not found, sessionPlayerKey=%s cmd=%d", sessionPlayerKey, cmd)
	}
	frame, err := clientSession.ProtocolCodec.Encode(cmd, index, body)
	if err != nil {
		return fmt.Errorf("encode transfer raw frame failed, sessionPlayerKey=%s cmd=%d err=%v", sessionPlayerKey, cmd, err)
	}
	if err := clientSession.SendRaw(frame); err != nil {
		return fmt.Errorf("send transfer raw response to client failed, sessionPlayerKey=%s cmd=%d err=%v", sessionPlayerKey, cmd, err)
	}
	logger.Info(fmt.Sprintf("send transfer raw response to client: cmd %d, bodyLen: %d", cmd, len(body)))
	return nil
}

func resolveBackendServerID(session *network.Session) int32 {
	if session == nil {
		return 0
	}
	if v, ok := session.GetAttr("serverId"); ok {
		return conv.Int32Value(v)
	}
	return 0
}

func newLogicIoDispatcher() network.IoDispatch {
	router := network.NewMessageRoute()
	ioDispatcher := &network.BaseIoDispatch{}
	ioDispatcher.AddHandler(&LogicRouter{router: router})
	return ioDispatcher
}
