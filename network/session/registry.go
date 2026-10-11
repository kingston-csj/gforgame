package session

import (
	"fmt"
	"sync"
)

var (
	sidMapper = sync.Map{}
	uidMap      = sync.Map{} // uid → sid
)

func BindSID(session Session) {
	sidMapper.Store(session.GetId(), session)
}

// Bind associates uid with an existing sid.
// Returns any previously bound session for the same uid (duplicate login).
func Bind(id ID, uid OwnerID) (Session, error) {
	if id == "" {
		return nil, fmt.Errorf("sid is empty")
	}
	if uid == "" {
		return nil, fmt.Errorf("uid is empty")
	}
	s, found := getSessionBySid(id)
	if !found {
		return nil, fmt.Errorf("session %s not found", id)
	}
	var oldSession Session
	if oldSID, found := getSid(uid); found && oldSID != id {
		oldSession, _ = getSessionBySid(oldSID)
	}
	uidMap.Store(uid, id)
	sidMapper.Store(id, s)
	return oldSession, nil
}

func Unbind(sid ID)  {
	s, found := popSession(sid)
	if !found {
		return
	}
	if nowSID, ok := getSid(s.GetOwnerId()); ok && nowSID == sid {
		uidMap.Delete(s.GetOwnerId())
	}
}

func popSession(sid ID) (Session, bool) {
	s, found := sidMapper.LoadAndDelete(sid)
	if !found {
		return nil, false
	}
	return s.(Session), true
}

func GetSessionByUid(uid OwnerID) (Session, bool) {
	sid, ok := getSid(uid)
	if !ok {
		return nil, false
	}
	return getSessionBySid(sid)
}

func getSessionBySid(id ID) (Session, bool) {
	session, ok := sidMapper.Load(id)
	if !ok {
		return nil, false
	}
	return session.(Session), ok
}

func getSid(uid OwnerID) (ID, bool) {
	sid, ok := uidMap.Load(uid)
	if !ok {
		return "", false
	}
	return sid.(ID), ok
}