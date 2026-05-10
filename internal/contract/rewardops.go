package contract

type RewardDefLite struct {
	Type  string
	Value string
}

type ItemRewardOps interface {
	AddByModelId(playerId string, itemId int32, amount int32) error
}

type CurrencyOps interface {
	Add(playerId string, kind string, amount int32)
}
