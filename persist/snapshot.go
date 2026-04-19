package persist

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Snapshotter allows entities to provide custom, safe snapshot logic.
// If not implemented, persist will fallback to generic JSON deep copy.
type Snapshotter interface {
	SnapshotEntity() (Entity, error)
}

func copyEntitySnapshot(entity Entity) (Entity, error) {
	if s, ok := entity.(Snapshotter); ok {
		return s.SnapshotEntity()
	}

	typ := reflect.TypeOf(entity)
	if typ == nil || typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("entity must be pointer to struct: %T", entity)
	}

	data, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	dst := reflect.New(typ.Elem()).Interface()
	if err = json.Unmarshal(data, dst); err != nil {
		return nil, err
	}

	cloned, ok := dst.(Entity)
	if !ok {
		return nil, fmt.Errorf("snapshot result is not Entity: %T", dst)
	}
	return cloned, nil
}
