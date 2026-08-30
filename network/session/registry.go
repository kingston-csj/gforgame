package session

import (
	"fmt"
	"sync"
)

var (
	sidMapper = sync.Map{}
	ownerId2IdMap = sync.Map{} // ownerId → sid
)

func BindSID(session Session) {
	sidMapper.Store(session.GetId(), session)
}

// Bind associates uid with an existing sid.
// Returns any previously bound session for the same uid (duplicate login).
func Bind(sid ID, ownerId OwnerID) (Session, error) {
	if sid == "" {
		return nil, fmt.Errorf("sid is empty")
	}
	if ownerId == "" {
		return nil, fmt.Errorf("uid is empty")
	}
	s, found := getSessionBySid(sid)
	if !found {
		return nil, fmt.Errorf("session %s not found", sid)
	}
	var oldSession Session
	if oldSID, found := getSessionId(ownerId); found && oldSID != sid {
		oldSession, _ = getSessionBySid(oldSID)
	}
	ownerId2IdMap.Store(ownerId, sid)
	sidMapper.Store(sid, s)
	return oldSession, nil
}

func Unbind(sid ID)  {
	s, found := popSession(sid)
	if !found {
		return
	}
	if nowSID, ok := getSessionId(s.GetOwnerId()); ok && nowSID == sid {
		ownerId2IdMap.Delete(s.GetOwnerId())
	}
}

func popSession(sid ID) (Session, bool) {
	s, found := sidMapper.LoadAndDelete(sid)
	if !found {
		return nil, false
	}
	return s.(Session), true
}

func GetSessionByOwnerId(ownerId OwnerID) (Session, bool) {
	sid, ok := getSessionId(ownerId)
	if !ok {
		return nil, false
	}
	return getSessionBySid(sid)
}

func getSessionBySid(sid ID) (Session, bool) {
	session, ok := sidMapper.Load(sid)
	if !ok {
		return nil, false
	}
	return session.(Session), ok
}

func getSessionId(ownerId OwnerID) (ID, bool) {
	sid, ok := ownerId2IdMap.Load(ownerId)
	if !ok {
		return "", false
	}
	return sid.(ID), ok
}