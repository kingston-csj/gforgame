package net

import (
	"sync"

	"github.com/forfun/gforgame/network/session"
)

const (
	AttrPlayerID = "player_id"
)

type OnlinePlayerRegistry struct {
	mu              sync.RWMutex
	player2Session  map[string]session.Session
	onlinePlayerIds map[string]bool
}

func NewOnlinePlayerRegistry() *OnlinePlayerRegistry {
	return &OnlinePlayerRegistry{
		mu:              sync.RWMutex{},
		player2Session:  make(map[string]session.Session),
		onlinePlayerIds: make(map[string]bool),
	}
}

func (r *OnlinePlayerRegistry) AddPlayerSession(session session.Session, playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	session.SetAttr(AttrPlayerID, playerID)
	r.player2Session[playerID] = session
}

// RemoveSession 移除session
// unbindPlayer 是否解绑玩家(如果是顶号，不应该解绑定)“
func (r *OnlinePlayerRegistry) RemovePlayerSession(session session.Session, unbindPlayer bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removePlayerSessionLocked(session, unbindPlayer)
}

func (r *OnlinePlayerRegistry) removePlayerSessionLocked(session session.Session, unbindPlayer bool) {
	playerID, ok := session.GetAttr(AttrPlayerID)
	if !ok {
		return
	}
	if unbindPlayer {
		delete(r.player2Session, playerID.(string))
	}
}

func (r *OnlinePlayerRegistry) GetSessionByPlayerID(playerID string) session.Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.player2Session[playerID]
}

func (r *OnlinePlayerRegistry) AddOnlinePlayer(playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.addOnlinePlayerUnsafe(playerID)
}

// 内部方法：约定 调用前必须已持有锁，不再加锁
func (r *OnlinePlayerRegistry) addOnlinePlayerUnsafe(playerID string) {
	r.onlinePlayerIds[playerID] = true
}

func (r *OnlinePlayerRegistry) RemoveOnlinePlayer(playerID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.onlinePlayerIds, playerID)
	return true
}

func (r *OnlinePlayerRegistry) IsOnline(playerID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.onlinePlayerIds[playerID]
}

func (r *OnlinePlayerRegistry) GetAllOnlinePlayerIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]string, 0, len(r.onlinePlayerIds))
	for id := range r.onlinePlayerIds {
		all = append(all, id)
	}
	return all
}
