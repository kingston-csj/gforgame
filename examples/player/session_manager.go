package player

import (
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/network"
)

var session_manager *SessionManager

var session2PlayerMap = make(map[*network.Session]*playerdomain.Player)
var player2SessionMap = make(map[*playerdomain.Player]*network.Session)

type SessionManager struct {
}

func GetSessionManager() *SessionManager {
	once.Do(func() {
		session_manager = &SessionManager{}
	})
	return session_manager
}

func (s *SessionManager) AddSession(session *network.Session, player *playerdomain.Player) {
	session2PlayerMap[session] = player
	player2SessionMap[player] = session
}

func (s *SessionManager) RemoveSession(session *network.Session) {
	player := session2PlayerMap[session]
	delete(session2PlayerMap, session)
	delete(player2SessionMap, player)
}

func (s *SessionManager) GetPlayerBySession(session *network.Session) *playerdomain.Player {
	return session2PlayerMap[session]
}

func (s *SessionManager) GetSessionByPlayer(player *playerdomain.Player) *network.Session {
	return player2SessionMap[player]
}
