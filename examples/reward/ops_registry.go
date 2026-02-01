package reward

import (
	itemcontract "io/github/gforgame/examples/contracts/item"
)

var (
    itemOps     itemcontract.ItemRewardOps
    currencyOps itemcontract.CurrencyOps
)

func SetItemOps(ops itemcontract.ItemRewardOps) {
    itemOps = ops
}

func SetSceneItemOps(ops itemcontract.ItemRewardOps) {
    itemOps = ops
}

func SetCurrencyOps(ops itemcontract.CurrencyOps) {
    currencyOps = ops
}

func getItemOps() itemcontract.ItemRewardOps {
    return itemOps
}

func getSceneItemOps() itemcontract.ItemRewardOps {
    return itemOps
}

func getCurrencyOps() itemcontract.CurrencyOps {
    return currencyOps
}
