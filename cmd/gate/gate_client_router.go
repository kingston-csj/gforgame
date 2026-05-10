package main

import (
	"fmt"

	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/jsonutil"
	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/internal/protos"
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
	if v, ok := session.GetAttr("id"); ok {
		if playerID, ok := v.(string); ok && playerID != "" {
			if current := network.GetSessionByPlayerId(playerID); current == session {
				unbindPlayerServer(playerID)
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
	if frame.Header.Cmd == protos.CmdReqPlayerLogin {
		var err error
		loginReq := &protos.ReqPlayerLogin{}
		if body, ok := frame.Msg.([]byte); ok {
			err = jsonutil.JsonBytesToStruct(body, loginReq)
		}
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

func HandleLoginReq(session *network.Session, loginReq *protos.ReqPlayerLogin, frame *protocol.RequestDataFrame) error {
	logger.Info(fmt.Sprintf("登录请求: %v", loginReq))
	if loginReq.PlayerId == "" {
		return fmt.Errorf("invalid login request: playerId is empty")
	}
	loginReq.ServerId = 1001
	if loginReq.ServerId <= 0 {
		return fmt.Errorf("invalid login request: serverId is empty")
	}
	if !isBackendServerConfigured(loginReq.ServerId) {
		return fmt.Errorf("target server not configured or not connected: serverId=%d", loginReq.ServerId)
	}

	oldSession := network.GetSessionByPlayerId(loginReq.PlayerId)
	if oldSession != nil && oldSession != session {
		logger.Info("玩家顶号登录[" + loginReq.PlayerId + "]")
		oldSession.SendAndClose(&protos.PushReplacingLogin{})
	}
	if oldSession == session {
		logger.Info("玩家重复登录[" + loginReq.PlayerId + "]")
	}
	session.SetAttr("id", loginReq.PlayerId)
	session.SetAttr("serverId", loginReq.ServerId)
	network.AddSession(session, loginReq.PlayerId)
	bindPlayerServer(loginReq.PlayerId, loginReq.ServerId)
	return transferMsgToLogic(session, frame, loginReq.ServerId)
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

	transfer := &protos.TransferGateToLogic{
		PlayerId: playerID,
		Cmd:      frame.Header.Cmd,
		Index:    frame.Header.Index,
		Body:     body,
	}
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

func isBackendServerConfigured(serverID int32) bool {
	serverType, ok := serverconfig.GetServerTypeByID(uint32(serverID))
	if !ok {
		return false
	}
	return serverType == logicServerType
}
