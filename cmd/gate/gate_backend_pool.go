package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/forfun/gforgame/common/logger"
	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/client"
)

// sendTransferToBackend 将网关转发消息发送到指定后端节点。
// 单连接模型下只尝试当前连接，失败后交给重连机制恢复。
func sendTransferToBackend(serverID int32, transfer any, index int32) error {
	session := pickBackendSession(serverID)
	if session == nil {
		// 没有可用连接时触发异步重连，消息由上层出站队列暂存重试。
		scheduleReconnect(serverID)
		return fmt.Errorf("backend session not ready, serverId=%d", serverID)
	}
	if err := session.Send(transfer, index); err != nil {
		// 发送失败通常意味着连接已失效，先摘除再重连。
		removeBackendSession(serverID, session)
		session.Close()
		scheduleReconnect(serverID)
		return err
	}
	return nil
}

// pickBackendSession 获取某个 serverID 当前可用的单连接。
// 如果连接已死亡会顺手清理为 nil。
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

// isSessionAlive 通过 Die 信号判断连接活性。
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

// removeBackendSession 仅在目标会话匹配时才清理，避免误删新连接。
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

// syncDiscoveredBackendPools 根据最新发现结果收敛连接池：
// 1. 新节点建立连接；
// 2. 已下线节点关闭连接；
// 3. 地址变化的节点切换到新地址重连。
func syncDiscoveredBackendPools(servers []DiscoveredServer) {
	desired := make(map[int32]string, len(servers))
	for _, server := range servers {
		if server.ID == 0 {
			continue
		}
		addr := strings.TrimSpace(buildDiscoveredServerAddr(server))
		if addr == "" {
			continue
		}
		desired[server.ID] = addr
	}
	reconcileBackendPools(desired)
}

// syncLocalBackendPools 仅根据本地配置同步后端连接。
// 该模式下只在启动阶段执行一次，不启用外部服务发现心跳。
func syncLocalBackendPools() {
	servers := serverconfig.GetServersByType(logicServerType)
	desired := make(map[int32]string, len(servers))
	for _, server := range servers {
		addr := strings.TrimSpace(server.Addr)
		if addr == "" {
			continue
		}
		desired[int32(server.ServerId)] = addr
	}
	reconcileBackendPools(desired)
}

func reconcileBackendPools(desired map[int32]string) {
	sessionsToClose := make([]*network.Session, 0)
	removedIDs := make([]int32, 0)
	changedIDs := make([]int32, 0)

	backendPoolsMu.Lock()
	for serverID, pool := range backendPools {
		addr, ok := desired[serverID]
		if !ok {
			if pool.session != nil {
				sessionsToClose = append(sessionsToClose, pool.session)
			}
			delete(backendPools, serverID)
			removedIDs = append(removedIDs, serverID)
			continue
		}
		if pool.addr != addr {
			pool.addr = addr
			pool.reconnecting = false
			if pool.session != nil {
				sessionsToClose = append(sessionsToClose, pool.session)
				pool.session = nil
			}
			changedIDs = append(changedIDs, serverID)
		}
	}
	backendPoolsMu.Unlock()

	for _, session := range sessionsToClose {
		if session != nil {
			session.Close()
		}
	}
	for _, serverID := range removedIDs {
		logger.Info(fmt.Sprintf("backend server removed: serverId=%d", serverID))
	}
	for _, serverID := range changedIDs {
		logger.Info(fmt.Sprintf("backend server address changed: serverId=%d", serverID))
	}
	for serverID, addr := range desired {
		ensureBackendPool(serverID, addr)
	}
}

// ensureBackendPool 确保节点状态存在并尽力建立单连接。
func ensureBackendPool(serverID int32, addr string) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return
	}
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
		// 首次连接失败也纳入统一重连路径。
		scheduleReconnect(serverID)
	}
}

// connectBackendSession 建立到后端节点的连接并替换旧连接。
// 单连接语义：若已有可用连接，新连接立即关闭。
func connectBackendSession(serverID int32, addr string) error {
	backendPoolsMu.RLock()
	pool := backendPools[serverID]
	if pool == nil {
		backendPoolsMu.RUnlock()
		return fmt.Errorf("backend pool removed, serverId=%d", serverID)
	}
	if strings.TrimSpace(pool.addr) != "" {
		addr = pool.addr
	}
	backendPoolsMu.RUnlock()

	socketClient := &client.WebSocketClient{RemoteAddress: addr, MsgCodec: gateMsgCodec}
	session, err := socketClient.OpenSession()
	if err != nil {
		return err
	}
	backendPoolsMu.Lock()
	pool = backendPools[serverID]
	if pool == nil {
		backendPoolsMu.Unlock()
		session.Close()
		return fmt.Errorf("backend pool removed, serverId=%d", serverID)
	}
	if pool.addr != addr {
		pool.reconnecting = false
		backendPoolsMu.Unlock()
		// 拨号期间地址已变化，直接丢弃旧连接结果。
		session.Close()
		return nil
	}
	if isSessionAlive(pool.session) {
		backendPoolsMu.Unlock()
		// 并发重连时可能已经有人先连上，直接放弃当前新连接。
		session.Close()
		return nil
	}
	oldSession := pool.session
	session.SetAttr("serverId", serverID)
	pool.session = session
	pool.reconnecting = false
	backendPoolsMu.Unlock()
	if oldSession != nil {
		// 替换成功后再关闭旧连接，避免短暂无连接窗口。
		oldSession.Close()
	}

	logger.Info(fmt.Sprintf("gate connected to backend server: serverId=%d addr=%s", serverID, addr))
	go consumeBackendSession(serverID, session)
	// 连接恢复后唤醒出站队列继续发送积压消息。
	notifyOutboundDispatcher()
	return nil
}

// consumeBackendSession 消费后端下行消息。
// 一旦连接断开，清理状态并触发延迟重连。
func consumeBackendSession(serverID int32, session *network.Session) {
	for {
		select {
		case <-session.Die:
			if removeBackendSession(serverID, session) {
				logger.Info(fmt.Sprintf("backend session closed and removed: serverId=%d", serverID))
				scheduleReconnect(serverID)
			}
			return
		case frame := <-session.DataReceived:
			if frame == nil {
				continue
			}
			if logicIoDispatcher != nil {
				logicIoDispatcher.OnMessageReceived(session, frame)
			}
		}
	}
}

// 关服时关闭所有后端连接。
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

// scheduleReconnect 基于事件驱动的延迟重连。
// 已被移除的节点不会再被重连。
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
			// 日志限流
			logMsg := fmt.Sprintf("reconnect backend failed, serverId=%d addr=%s err=%v", sid, remoteAddr, err)
			if !logFilter.IsExists(logMsg) {
				logger.ErrorNoStack(logMsg)
			}

			backendPoolsMu.Lock()
			if p := backendPools[sid]; p != nil {
				p.reconnecting = false
			}
			backendPoolsMu.Unlock()
			scheduleReconnect(sid)
		}
	}(serverID, addr)
}

// startBackendSessionMonitor 定时巡检所有后端连接。
// 若发现连接已断开，会立刻走统一的重连流程恢复。
func startBackendSessionMonitor() {
	go func() {
		ticker := time.NewTicker(time.Duration(backendMonitorInterval) * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				checkBackendSessions()
			case <-backendMonitorStop:
				return
			}
		}
	}()
}

func checkBackendSessions() {
	serverIDs := make([]int32, 0)
	sessionsToClose := make([]*network.Session, 0)

	backendPoolsMu.Lock()
	for serverID, pool := range backendPools {
		if pool == nil {
			continue
		}
		if isSessionAlive(pool.session) {
			continue
		}
		if pool.session != nil {
			sessionsToClose = append(sessionsToClose, pool.session)
			pool.session = nil
		}
		serverIDs = append(serverIDs, serverID)
	}
	backendPoolsMu.Unlock()

	for _, session := range sessionsToClose {
		if session != nil {
			session.Close()
		}
	}
	for _, serverID := range serverIDs {
		scheduleReconnect(serverID)
	}
}
