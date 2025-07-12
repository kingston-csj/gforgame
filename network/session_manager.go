package network

import (
	"sync"

	"io/github/gforgame/domain"
)

var (
	sessionManager *SessionManager
	once           sync.Once
)

var (
	session2PlayerMap = make(map[*Session]domain.Player)
	player2SessionMap = make(map[string]*Session)
)

type SessionManager struct{}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

func GetSessionManager() *SessionManager {
	once.Do(func() {
		sessionManager = NewSessionManager()
	})
	return sessionManager
}

func AddSession(session *Session, player domain.Player) {
	session2PlayerMap[session] = player
	player2SessionMap[player.GetId()] = session
}

func RemoveSession(session *Session) {
	player := session2PlayerMap[session]
	delete(session2PlayerMap, session)
	delete(player2SessionMap, player.GetId())
}

func GetPlayerBySession(session *Session) domain.Player {
	return session2PlayerMap[session]
}

func GetSessionByPlayerId(playerId string) *Session {
	return player2SessionMap[playerId]
}

func GetAllSessions() []*Session {
	allSessions := make([]*Session, 0, len(session2PlayerMap))
	for session := range session2PlayerMap {
		allSessions = append(allSessions, session)
	}
	return allSessions
}

func GetAllOnlinePlayers() []domain.Player {
	allPlayers := make([]domain.Player, 0, len(session2PlayerMap))
	for _, player := range session2PlayerMap {
		allPlayers = append(allPlayers, player)
	}
	return allPlayers
}

func GetAllOnlinePlayerSessions() []*Session {
	allSessions := make([]*Session, 0, len(session2PlayerMap))
	for session := range session2PlayerMap {
		allSessions = append(allSessions, session)
	}
	return allSessions
}

func IsOnline(playerId string) bool {
	return player2SessionMap[playerId] != nil
}

