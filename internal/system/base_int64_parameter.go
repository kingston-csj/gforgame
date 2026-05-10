package system

import (
	"strconv"
	"sync"
	"sync/atomic"
)

type baseInt64Parameter struct {
	ID       string       `json:"id"`
	value    atomic.Int64 `json:"-"`
	loadOnce sync.Once    `json:"-"`
}

func (b *baseInt64Parameter) init(id string) {
	if b.ID == "" {
		b.ID = id
	}
}

func (b *baseInt64Parameter) getID() string {
	return b.ID
}

func (b *baseInt64Parameter) getValue(loadFn func() string) int64 {
	b.loadOnce.Do(func() {
		raw := loadFn()
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

func (b *baseInt64Parameter) parseFromStore(loadFn func() string) int64 {
	raw := loadFn()
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

func (b *baseInt64Parameter) saveValue(v int64, payload any) {
	b.value.Store(v)
	saveSystemParameterValue(b.getID(), strconv.FormatInt(v, 10), payload)
}
