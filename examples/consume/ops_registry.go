package consume

import (
	"sync"

	itemcontract "github.com/forfun/gforgame/examples/contract"
)

var (
    mu          sync.RWMutex
    itemOps     itemcontract.ItemConsumeOps
)

func SetItemOps(ops itemcontract.ItemConsumeOps) {
    mu.Lock()
    defer mu.Unlock()
    itemOps = ops
}

func GetItemOps() itemcontract.ItemConsumeOps {
    mu.RLock()
    defer mu.RUnlock()
    return itemOps
}