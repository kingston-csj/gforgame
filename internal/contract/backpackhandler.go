package contract

type BackpackHandler interface {

	// 通过道具模型id进行扣除（后端接口）
	UseByModelId(playerId string, itemId int32, count int32) error

	// 通过道具uid进行扣除（前端接口）
	UseByUid(playerId string, itemUid string, count int32) (error, []RewardDefLite)

	AddByModelId(playerId string, itemId int32, count int32) error
}