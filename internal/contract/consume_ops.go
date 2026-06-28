package contract

type ItemConsumeOps interface {
	UseByModelId(playerId string, itemId int32, amount int32) error
}
