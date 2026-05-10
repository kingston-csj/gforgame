package system

import (
	"sync"
	"sync/atomic"
)

type baseStringParameter struct {
	ID       string       `json:"id"`
	value    atomic.Value `json:"-"`
	loadOnce sync.Once    `json:"-"`
}

func (b *baseStringParameter) init(id string) {
	if b.ID == "" {
		b.ID = id
	}
}

func (b *baseStringParameter) getID() string {
	return b.ID
}

func (b *baseStringParameter) getValue(loadFn func() string) string {
	b.loadOnce.Do(func() {
		raw := loadFn()
		if raw == "" {
			b.value.Store("")
			return
		}
		b.value.Store(raw)
	})
	return b.value.Load().(string)
}

func (b *baseStringParameter) parseFromStore(loadFn func() string) string {
	raw := loadFn()
	if raw == "" {
		return ""
	}
	b.value.Store(raw)
	return raw
}

// saveValue 保存值并持久化
func (b *baseStringParameter) saveValue(v string, payload any) {
	b.value.Store(v)
	saveSystemParameterValue(b.getID(), v, payload)
}
