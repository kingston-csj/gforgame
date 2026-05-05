package buff

import (
	"github.com/forfun/gforgame/common/util"
	"github.com/forfun/gforgame/examples/config"
	configdomain "github.com/forfun/gforgame/examples/domain/config"
	"github.com/forfun/gforgame/examples/fight/state"
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
