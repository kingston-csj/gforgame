package buff

import (
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/fight/state"
	"io/github/gforgame/util"
)

type StateBuff struct {
	Buff
	StateType state.StateType
}

func NewStateBuff(modelId int32) *StateBuff {
	buffData := context.GetConfigRecordAs[config.BuffData]("buff", int64(modelId))
	return &StateBuff{
		Buff: Buff{
			ModelId:  modelId,
			Id:       util.GetNextId(),
			Duration: buffData.Duration,
		},
		StateType: state.StateType(buffData.Params),
	}
}
