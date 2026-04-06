package network

import (
	"net"
	"sync"
)

// sessionRegistry 统一维护连接、会话、玩家三者关系（包内私有）
type sessionRegistry struct {
	mu sync.RWMutex

	conn2Session   map[net.Conn]*Session
	session2Player map[*Session]string
	player2Session map[string]*Session
}

func newSessionRegistry() *sessionRegistry {
	return &sessionRegistry{
		conn2Session:   make(map[net.Conn]*Session),
		session2Player: make(map[*Session]string),
		player2Session: make(map[string]*Session),
	}
}

var globalSessionRegistry = newSessionRegistry()

func (r *sessionRegistry) registerConnSession(conn net.Conn, s *Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conn2Session[conn] = s
}

func (r *sessionRegistry) getSessionByConn(conn net.Conn) *Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.conn2Session[conn]
}

func (r *sessionRegistry) unregisterConn(conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	s := r.conn2Session[conn]
	delete(r.conn2Session, conn)
	if s == nil {
		return
	}
	r.removePlayerSessionLocked(s)
}

func (r *sessionRegistry) addPlayerSession(session *Session, playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.session2Player[session] = playerID
	r.player2Session[playerID] = session
}

func (r *sessionRegistry) removePlayerSession(session *Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removePlayerSessionLocked(session)
}

func (r *sessionRegistry) removePlayerSessionLocked(session *Session) {
	playerID := r.session2Player[session]
	delete(r.session2Player, session)
	if playerID == "" {
		return
	}
	delete(r.player2Session, playerID)
}

func (r *sessionRegistry) getPlayerIDBySession(session *Session) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	playerID, ok := r.session2Player[session]
	return playerID, ok
}

func (r *sessionRegistry) getSessionByPlayerID(playerID string) *Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.player2Session[playerID]
}

func (r *sessionRegistry) getAllPlayerSessions() []*Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]*Session, 0, len(r.session2Player))
	for s := range r.session2Player {
		all = append(all, s)
	}
	return all
}

func (r *sessionRegistry) getAllOnlinePlayerIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]string, 0, len(r.player2Session))
	for id := range r.player2Session {
		all = append(all, id)
	}
	return all
}

func (r *sessionRegistry) isOnline(playerID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.player2Session[playerID] != nil
}

func (r *sessionRegistry) getAllConnSessions() []*Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sessions := make([]*Session, 0, len(r.conn2Session))
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

func RegisterSession(conn net.Conn, s *Session) {
	globalSessionRegistry.registerConnSession(conn, s)
}

func GetSession(conn net.Conn) *Session {
	return globalSessionRegistry.getSessionByConn(conn)
}

func UnregisterSession(conn net.Conn) {
	globalSessionRegistry.unregisterConn(conn)
}

func CloseAllSessions() {
	globalSessionRegistry.closeAllSessions()
}

func AddSession(session *Session, playerID string) {
	globalSessionRegistry.addPlayerSession(session, playerID)
}

func RemoveSession(session *Session) {
	globalSessionRegistry.removePlayerSession(session)
}

func GetPlayerIDBySession(session *Session) (string, bool) {
	return globalSessionRegistry.getPlayerIDBySession(session)
}

func GetSessionByPlayerId(playerID string) *Session {
	return globalSessionRegistry.getSessionByPlayerID(playerID)
}

func GetAllSessions() []*Session {
	return globalSessionRegistry.getAllPlayerSessions()
}

func GetAllOnlinePlayerIds() []string {
	return globalSessionRegistry.getAllOnlinePlayerIDs()
}

func GetAllOnlinePlayerSessions() []*Session {
	return globalSessionRegistry.getAllPlayerSessions()
}

func IsOnline(playerID string) bool {
	return globalSessionRegistry.isOnline(playerID)
}
