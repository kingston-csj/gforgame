package session

import (
	"net"
	"sync"
)

// sessionRegistry 统一维护连接、会话、玩家三者关系（包内私有）。
// 内部以 *BaseSession 作为身份标识，外部暴露 Session 接口。
type sessionRegistry struct {
	mu sync.RWMutex

	conn2Session    map[net.Conn]*BaseSession
	session2Player  map[*BaseSession]string
	player2Session  map[string]*BaseSession
	onlinePlayerIds map[string]bool
}

func newSessionRegistry() *sessionRegistry {
	return &sessionRegistry{
		conn2Session:    make(map[net.Conn]*BaseSession),
		session2Player:  make(map[*BaseSession]string),
		player2Session:  make(map[string]*BaseSession),
		onlinePlayerIds: make(map[string]bool),
	}
}

var globalSessionRegistry = newSessionRegistry()

// toBase 同包内统一把接口转成具体实现，不存在则 panic（编程错误）。
func toBase(session Session) *BaseSession {
	base, ok := session.(*BaseSession)
	if !ok {
		panic("session: unexpected Session implementation, only *BaseSession is supported")
	}
	return base
}

func (r *sessionRegistry) registerConnSession(conn net.Conn, s *BaseSession) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conn2Session[conn] = s
}

func (r *sessionRegistry) getSessionByConn(conn net.Conn) *BaseSession {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.conn2Session[conn]
}

func (r *sessionRegistry) unregisterConn(conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.conn2Session, conn)
}

func (r *sessionRegistry) addPlayerSession(session *BaseSession, playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.session2Player[session] = playerID
	r.player2Session[playerID] = session
	r.addOnlinePlayerUnsafe(playerID)
}

func (r *sessionRegistry) addOnlinePlayer(playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.addOnlinePlayerUnsafe(playerID)
}

// 内部方法：约定 调用前必须已持有锁，不再加锁
func (r *sessionRegistry) addOnlinePlayerUnsafe(playerID string) {
	r.onlinePlayerIds[playerID] = true
}
func (r *sessionRegistry) removePlayerSession(session *BaseSession, unbindPlayer bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removePlayerSessionLocked(session, unbindPlayer)
}

func (r *sessionRegistry) removePlayerSessionLocked(session *BaseSession, unbindPlayer bool) {
	playerID := r.session2Player[session]
	delete(r.session2Player, session)
	if playerID == "" {
		return
	}
	if unbindPlayer {
		delete(r.player2Session, playerID)
		delete(r.onlinePlayerIds, playerID)
	}
}

func (r *sessionRegistry) removeOnlinePlayer(playerID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.onlinePlayerIds, playerID)
	return true
}

func (r *sessionRegistry) getPlayerIDBySession(session *BaseSession) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	playerID, ok := r.session2Player[session]
	return playerID, ok
}

func (r *sessionRegistry) getSessionByPlayerID(playerID string) *BaseSession {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.player2Session[playerID]
}

func (r *sessionRegistry) getAllPlayerSessions() []*BaseSession {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]*BaseSession, 0, len(r.session2Player))
	for s := range r.session2Player {
		all = append(all, s)
	}
	return all
}

func (r *sessionRegistry) getAllOnlinePlayerIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]string, 0, len(r.onlinePlayerIds))
	for id := range r.onlinePlayerIds {
		all = append(all, id)
	}
	return all
}

func (r *sessionRegistry) isOnline(playerID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.onlinePlayerIds[playerID]
}

func (r *sessionRegistry) getAllConnSessions() []*BaseSession {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sessions := make([]*BaseSession, 0, len(r.conn2Session))
	for _, s := range r.conn2Session {
		sessions = append(sessions, s)
	}
	return sessions
}

func (r *sessionRegistry) closeAllSessions() {
	for _, s := range r.getAllConnSessions() {
		s.Close()
	}
}

// -------- 对外导出函数：参数/返回值均使用 Session 接口，边界处转实现 --------

func RegisterSession(conn net.Conn, s Session) {
	globalSessionRegistry.registerConnSession(conn, toBase(s))
}

func GetSession(conn net.Conn) Session {
	s := globalSessionRegistry.getSessionByConn(conn)
	if s == nil {
		return nil
	}
	return s
}

func UnregisterSession(conn net.Conn) {
	globalSessionRegistry.unregisterConn(conn)
}

func CloseAllSessions() {
	globalSessionRegistry.closeAllSessions()
}

func AddSession(session Session, playerID string) {
	globalSessionRegistry.addPlayerSession(toBase(session), playerID)
}

func AddOnlinePlayer(playerID string) {
	globalSessionRegistry.addOnlinePlayer(playerID)
}

// RemoveSession 移除session
// unbindPlayer 是否解绑玩家(如果是顶号，不应该解绑定)
func RemoveSession(session Session, unbindPlayer bool) {
	globalSessionRegistry.removePlayerSession(toBase(session), unbindPlayer)
}

func GetPlayerIDBySession(session Session) (string, bool) {
	return globalSessionRegistry.getPlayerIDBySession(toBase(session))
}

func GetSessionByPlayerId(playerID string) Session {
	s := globalSessionRegistry.getSessionByPlayerID(playerID)
	if s == nil {
		return nil
	}
	return s
}

func GetAllSessions() []Session {
	raw := globalSessionRegistry.getAllPlayerSessions()
	all := make([]Session, 0, len(raw))
	for _, s := range raw {
		all = append(all, s)
	}
	return all
}

func GetAllOnlinePlayerIds() []string {
	return globalSessionRegistry.getAllOnlinePlayerIDs()
}

func GetAllOnlinePlayerSessions() []Session {
	raw := globalSessionRegistry.getAllPlayerSessions()
	all := make([]Session, 0, len(raw))
	for _, s := range raw {
		all = append(all, s)
	}
	return all
}

func IsOnline(playerID string) bool {
	return globalSessionRegistry.isOnline(playerID)
}

func RemoveOnlinePlayer(playerID string) bool {
	return globalSessionRegistry.removeOnlinePlayer(playerID)
}
