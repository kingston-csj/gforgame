package system

import (
	"sync"
	"sync/atomic"

	systemrepo "github.com/forfun/gforgame/internal/infra/repository/system"
)

type baseStringParameter struct {
	ID         string                       `json:"id"`
	repository *systemrepo.SystemRepository `json:"-"`
	value      atomic.Value                 `json:"-"`
	loadOnce   sync.Once                    `json:"-"`
}

func (b *baseStringParameter) init(id string, repo *systemrepo.SystemRepository) {
	if b.ID == "" {
		b.ID = id
	}
	if b.repository == nil {
		b.repository = repo
	}
}

func (b *baseStringParameter) getID() string {
	return b.ID
}

func (b *baseStringParameter) getValue() string {
	b.loadOnce.Do(func() {
		raw := b.loadValue()
		if raw == "" {
			b.value.Store("")
			return
		}
		b.value.Store(raw)
	})
	return b.value.Load().(string)
}

func (b *baseStringParameter) parseFromStore() string {
	raw := b.loadValue()
	if raw == "" {
		return ""
	}
	b.value.Store(raw)
	return raw
}

// saveValue 保存值并持久化
func (b *baseStringParameter) saveValue(v string) {
	b.value.Store(v)
	if b.repository == nil {
		return
	}
	b.repository.SaveValue(b.getID(), v)
}

func (b *baseStringParameter) loadValue() string {
	if b.repository == nil {
		return ""
	}
	return b.repository.LoadValue(b.getID())
}
