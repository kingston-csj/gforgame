package system

import (
	"strconv"
	"sync"
	"sync/atomic"

	systemrepo "github.com/forfun/gforgame/internal/infra/repository/system"
)

type baseInt64Parameter struct {
	ID         string                       `json:"id"`
	repository *systemrepo.SystemRepository `json:"-"`
	value      atomic.Int64                 `json:"-"`
	loadOnce   sync.Once                    `json:"-"`
}

func (b *baseInt64Parameter) init(id string, repo *systemrepo.SystemRepository) {
	if b.ID == "" {
		b.ID = id
	}
	if b.repository == nil {
		b.repository = repo
	}
}

func (b *baseInt64Parameter) getID() string {
	return b.ID
}

func (b *baseInt64Parameter) getValue() int64 {
	b.loadOnce.Do(func() {
		raw := b.loadValue()
		if raw == "" {
			b.value.Store(0)
			return
		}
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			b.value.Store(0)
			return
		}
		b.value.Store(v)
	})
	return b.value.Load()
}

func (b *baseInt64Parameter) parseFromStore() int64 {
	raw := b.loadValue()
	if raw == "" {
		return 0
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0
	}
	b.value.Store(v)
	return v
}

func (b *baseInt64Parameter) saveValue(v int64) {
	b.value.Store(v)
	if b.repository == nil {
		return
	}
	b.repository.SaveValue(b.getID(), strconv.FormatInt(v, 10))
}

func (b *baseInt64Parameter) loadValue() string {
	if b.repository == nil {
		return ""
	}
	return b.repository.LoadValue(b.getID())
}
