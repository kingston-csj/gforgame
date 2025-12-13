package reward

import (
	itemcontract "io/github/gforgame/examples/contracts/item"
	"sync"
)

var (
    mu          sync.RWMutex
    itemOps     itemcontract.ItemOps
    currencyOps itemcontract.CurrencyOps
)

func SetItemOps(ops itemcontract.ItemOps) {
    mu.Lock()
    itemOps = ops
    mu.Unlock()
}

func SetCurrencyOps(ops itemcontract.CurrencyOps) {
    mu.Lock()
    currencyOps = ops
    mu.Unlock()
}

func getItemOps() itemcontract.ItemOps {
    mu.RLock()
    defer mu.RUnlock()
    return itemOps
}

func getCurrencyOps() itemcontract.CurrencyOps {
    mu.RLock()
    defer mu.RUnlock()
    return currencyOps
}
