package session

import (
	"io/github/gforgame/examples/types"
	"io/github/gforgame/network"
	"sync"
)

var (
	sessionManager *SessionManager
	once           sync.Once
)

var (
	session2PlayerMap = make(map[*network.Session]types.Player)
	player2SessionMap = make(map[string]*network.Session)
)

type SessionManager struct {
}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

func GetSessionManager() *SessionManager {
	once.Do(func() {
		sessionManager = NewSessionManager()
	})
	return sessionManager
}

func (s *SessionManager) AddSession(session *network.Session, player types.Player) {
	session2PlayerMap[session] = player
	player2SessionMap[player.GetID()] = session
}

func (s *SessionManager) RemoveSession(session *network.Session) {
	player := session2PlayerMap[session]
	delete(session2PlayerMap, session)
	delete(player2SessionMap, player.GetID())
}

func (s *SessionManager) GetPlayerBySession(session *network.Session) types.Player {
	return session2PlayerMap[session]
}

func (s *SessionManager) GetSessionByPlayerId(playerId string) *network.Session {
	return player2SessionMap[playerId]
}
