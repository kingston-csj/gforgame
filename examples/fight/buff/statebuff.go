package buff

import (
	"io/github/gforgame/examples/config"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/fight/state"
	"io/github/gforgame/util"
)

type StateBuff struct {
	Buff
	StateType state.StateType
}

func NewStateBuff(modelId int32) *StateBuff {
	buffData := config.QueryById[configdomain.BuffData](modelId)
	return &StateBuff{
		Buff: Buff{
			ModelId:  modelId,
			Id:       util.GetNextID(),
			Duration: buffData.Duration,
		},
		StateType: state.StateType(buffData.Params),
	}
}
