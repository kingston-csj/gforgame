package reward

import (
    "io/github/gforgame/examples/domain/contracts"
    "sync"
)

var (
    mu          sync.RWMutex
    itemOps     contracts.ItemOps
    currencyOps contracts.CurrencyOps
)

func SetItemOps(ops contracts.ItemOps) {
    mu.Lock()
    itemOps = ops
    mu.Unlock()
}

func SetCurrencyOps(ops contracts.CurrencyOps) {
    mu.Lock()
    currencyOps = ops
    mu.Unlock()
}

func getItemOps() contracts.ItemOps {
    mu.RLock()
    defer mu.RUnlock()
    return itemOps
}

func getCurrencyOps() contracts.CurrencyOps {
    mu.RLock()
    defer mu.RUnlock()
    return currencyOps
}
