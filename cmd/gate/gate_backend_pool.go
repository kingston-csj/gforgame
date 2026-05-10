package main

import (
	"fmt"
	"time"

	"github.com/forfun/gforgame/common/logger"
	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/client"
)

func sendTransferToBackend(serverID int32, transfer *protos.TransferGateToLogic, index int32) error {
	session := pickBackendSession(serverID)
	if session == nil {
		scheduleReconnect(serverID)
		return fmt.Errorf("backend session not ready, serverId=%d", serverID)
	}
	if err := session.Send(transfer, index); err != nil {
		removeBackendSession(serverID, session)
		session.Close()
		scheduleReconnect(serverID)
		return err
	}
	return nil
}

func pickBackendSession(serverID int32) *network.Session {
	backendPoolsMu.Lock()
	defer backendPoolsMu.Unlock()
	pool, ok := backendPools[serverID]
	if !ok {
		return nil
	}
	if !isSessionAlive(pool.session) {
		pool.session = nil
	}
	return pool.session
}

func isSessionAlive(session *network.Session) bool {
	if session == nil {
		return false
	}
	select {
	case <-session.Die:
		return false
	default:
		return true
	}
}

func removeBackendSession(serverID int32, target *network.Session) bool {
	backendPoolsMu.Lock()
	defer backendPoolsMu.Unlock()
	pool := backendPools[serverID]
	if pool == nil || pool.session == nil {
		return false
	}
	if target != nil && pool.session != target {
		return false
	}
	pool.session = nil
	return true
}

func startLogicConnectScheduler() {
	registerBackendServers()
}

func registerBackendServers() {
	servers := serverconfig.GetServersByType(logicServerType)
	for _, server := range servers {
		serverID := int32(server.ServerId)
		ensureBackendPool(serverID, server.Addr)
	}
}

func ensureBackendPool(serverID int32, addr string) {
	backendPoolsMu.Lock()
	pool := backendPools[serverID]
	if pool == nil {
		pool = &backendPool{serverID: serverID, addr: addr}
		backendPools[serverID] = pool
	} else {
		pool.addr = addr
	}
	connected := isSessionAlive(pool.session)
	backendPoolsMu.Unlock()
	if connected {
		return
	}
	if err := connectBackendSession(serverID, addr); err != nil {
		logger.ErrorNoStack(fmt.Errorf("connect backend failed, serverId=%d addr=%s err=%v", serverID, addr, err))
		scheduleReconnect(serverID)
	}
}

func connectBackendSession(serverID int32, addr string) error {
	tcpClient := &client.TcpSocketClient{RemoteAddress: addr, MsgCodec: gateMsgCodec}
	session, err := tcpClient.OpenSession()
	if err != nil {
		return err
	}
	backendPoolsMu.Lock()
	pool := backendPools[serverID]
	if pool == nil {
		pool = &backendPool{serverID: serverID, addr: addr}
		backendPools[serverID] = pool
	}
	if isSessionAlive(pool.session) {
		backendPoolsMu.Unlock()
		session.Close()
		return nil
	}
	oldSession := pool.session
	pool.session = session
	pool.reconnecting = false
	backendPoolsMu.Unlock()
	if oldSession != nil {
		oldSession.Close()
	}

	logger.Info(fmt.Sprintf("gate connected to backend server: serverId=%d addr=%s", serverID, addr))
	go consumeBackendSession(serverID, session)
	notifyOutboundDispatcher()
	return nil
}

func consumeBackendSession(serverID int32, session *network.Session) {
	for frame := range session.DataReceived {
		if logicIoDispatcher != nil {
			logicIoDispatcher.OnMessageReceived(session, frame)
		}
	}
	if removeBackendSession(serverID, session) {
		logger.Info(fmt.Sprintf("backend session closed and removed: serverId=%d", serverID))
		scheduleReconnect(serverID)
	}
}

func closeAllBackendSessions() {
	backendPoolsMu.Lock()
	all := make([]*network.Session, 0)
	for _, pool := range backendPools {
		if pool.session != nil {
			all = append(all, pool.session)
		}
	}
	backendPools = make(map[int32]*backendPool)
	backendPoolsMu.Unlock()
	for _, session := range all {
		if session != nil {
			session.Close()
		}
	}
}

func scheduleReconnect(serverID int32) {
	backendPoolsMu.Lock()
	pool := backendPools[serverID]
	if pool == nil {
		backendPoolsMu.Unlock()
		return
	}
	if pool.reconnecting {
		backendPoolsMu.Unlock()
		return
	}
	pool.reconnecting = true
	addr := pool.addr
	backendPoolsMu.Unlock()

	go func(sid int32, remoteAddr string) {
		time.Sleep(time.Duration(backendReconnectDelay) * time.Millisecond)
		if err := connectBackendSession(sid, remoteAddr); err != nil {
			logger.ErrorNoStack(fmt.Errorf("reconnect backend failed, serverId=%d addr=%s err=%v", sid, remoteAddr, err))
			backendPoolsMu.Lock()
			if p := backendPools[sid]; p != nil {
				p.reconnecting = false
			}
			backendPoolsMu.Unlock()
			scheduleReconnect(sid)
		}
	}(serverID, addr)
}
