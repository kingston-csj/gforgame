package consume

import (
	itemcontract "io/github/gforgame/examples/contracts/item"
	"sync"
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