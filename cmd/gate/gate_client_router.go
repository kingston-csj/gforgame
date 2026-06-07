package main

import (
	"fmt"

	"github.com/forfun/gforgame/common/logger"
	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/gateway/contract"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
)

// 作为网关，处理客户端请求
type ClientRouter struct {
	router *network.MessageRoute
}

type MyMessageDispatch struct {
	network.BaseIoDispatch
}

func (m *MyMessageDispatch) OnSessionClosed(session *network.Session) {
	if v, ok := session.GetAttr("sessionPlayerKey"); ok {
		if sessionPlayerKey, ok := v.(string); ok && sessionPlayerKey != "" {
			if current := network.GetSessionByPlayerId(sessionPlayerKey); current == session {
				unbindPlayerServer(sessionPlayerKey)
			}
		}
	}
	network.RemoveSession(session)
}

func (g *ClientRouter) MessageReceived(session *network.Session, frame *protocol.RequestDataFrame) bool {
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorNoStack(fmt.Errorf("panic recovered: %v", r))
		}
	}()
	msgName, _ := network.GetMsgName(frame.Header.Cmd)
	logger.Info(fmt.Sprintf("接收消息: cmd:%d, name:%s, 内容:%s", frame.Header.Cmd, msgName, frame.Msg))
	if frame.Header.Cmd == gateLoginAdapter.LoginCmd() {
		body, ok := frame.Msg.([]byte)
		if !ok {
			logger.ErrorNoStack(fmt.Errorf("login payload type invalid: %T", frame.Msg))
			return false
		}
		loginReq, err := gateLoginAdapter.DecodeLoginRequest(body)
		if err != nil {
			logger.ErrorNoStack(fmt.Errorf("unmarshal login req failed: %v", err))
			return false
		}
		logger.Info(fmt.Sprintf("登录请求: %v", loginReq))
		if err := HandleLoginReq(session, loginReq, frame); err != nil {
			logger.ErrorNoStack(err)
		}
		return true
	}
	if err := transferMsgToLogic(session, frame, 0); err != nil {
		logger.ErrorNoStack(err)
	}
	return true
}

func HandleLoginReq(session *network.Session, loginReq contract.GateLoginRequest, frame *protocol.RequestDataFrame) error {
	logger.Info(fmt.Sprintf("登录请求: %v", loginReq))
	playerID := loginReq.GetPlayerID()
	serverID := loginReq.GetServerID()
	if playerID == "" {
		return fmt.Errorf("invalid login request: playerId is empty")
	}
	if serverID <= 0 {
		return fmt.Errorf("invalid login request: serverId is empty")
	}
	if !isBackendServerConfigured(serverID) {
		return fmt.Errorf("target server not configured or not connected: serverId=%d", serverID)
	}

	sessionPlayerKey := buildSessionPlayerKey(serverID, playerID)
	oldSession := network.GetSessionByPlayerId(sessionPlayerKey)
	if oldSession != nil && oldSession != session {
		logger.Info("玩家顶号登录[" + sessionPlayerKey + "]")
		oldSession.SendAndClose(gateLoginAdapter.NewReplacingLoginPush())
	}
	if oldSession == session {
		logger.Info("玩家重复登录[" + sessionPlayerKey + "]")
	}
	session.SetAttr("id", playerID)
	session.SetAttr("serverId", serverID)
	session.SetAttr("sessionPlayerKey", sessionPlayerKey)
	network.AddSession(session, sessionPlayerKey)
	bindPlayerServer(sessionPlayerKey, serverID)
	return transferMsgToLogic(session, frame, serverID)
}

// 转发消息到logic层
func transferMsgToLogic(session *network.Session, frame *protocol.RequestDataFrame, forceServerID int32) error {
	body, ok := frame.Msg.([]byte)
	if !ok {
		return fmt.Errorf("transfer payload type invalid, cmd=%d type=%T", frame.Header.Cmd, frame.Msg)
	}

	playerID := ""
	if v, ok := session.GetAttr("id"); ok {
		if s, ok := v.(string); ok {
			playerID = s
		}
	}

	serverID := forceServerID
	if serverID <= 0 {
		serverID = resolvePlayerServerID(session, playerID)
	}
	if serverID <= 0 {
		return fmt.Errorf("target server is empty, playerId=%s cmd=%d", playerID, frame.Header.Cmd)
	}

	transfer := gateTransferCodec.NewTransferMessage(playerID, frame.Header.Cmd, frame.Header.Index, body)
	if err := enqueueTransfer(serverID, transfer, frame.Header.Index); err != nil {
		return fmt.Errorf("forward message to backend failed, serverId=%d cmd=%d err=%v", serverID, frame.Header.Cmd, err)
	}
	return nil
}

func bindPlayerServer(playerID string, serverID int32) {
	playerServerIDMapMu.Lock()
	defer playerServerIDMapMu.Unlock()
	playerServerIDMap[playerID] = serverID
}

func unbindPlayerServer(playerID string) {
	playerServerIDMapMu.Lock()
	defer playerServerIDMapMu.Unlock()
	delete(playerServerIDMap, playerID)
}

func getPlayerServer(playerID string) (int32, bool) {
	playerServerIDMapMu.RLock()
	defer playerServerIDMapMu.RUnlock()
	serverID, ok := playerServerIDMap[playerID]
	return serverID, ok
}

func resolvePlayerServerID(session *network.Session, playerID string) int32 {
	if v, ok := session.GetAttr("sessionPlayerKey"); ok {
		if sessionPlayerKey, ok := v.(string); ok && sessionPlayerKey != "" {
			if serverID, ok := getPlayerServer(sessionPlayerKey); ok {
				return serverID
			}
		}
	}
	if playerID != "" {
		if serverID, ok := getPlayerServer(playerID); ok {
			return serverID
		}
	}
	if v, ok := session.GetAttr("serverId"); ok {
		switch sid := v.(type) {
		case int32:
			return sid
		case int:
			return int32(sid)
		}
	}
	return 0
}

func buildSessionPlayerKey(serverID int32, playerID string) string {
	return fmt.Sprintf("%d_%s", serverID, playerID)
}

func isBackendServerConfigured(serverID int32) bool {
	serverType, ok := serverconfig.GetServerTypeByID(uint32(serverID))
	if !ok {
		return false
	}
	return serverType == logicServerType
}
