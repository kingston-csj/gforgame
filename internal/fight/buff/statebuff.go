package buff

import (
	"github.com/forfun/gforgame/internal/config"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	"github.com/forfun/gforgame/internal/fight/state"
	"github.com/forfun/gforgame/internal/idgen"
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
			Id:       idgen.GetNextID(),
			Duration: buffData.Duration,
		},
		StateType: state.StateType(buffData.Params),
	}
}
