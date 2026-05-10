package buff

import (
	"strings"

	"github.com/forfun/gforgame/common/util"
	"github.com/forfun/gforgame/internal/config"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	"github.com/forfun/gforgame/internal/fight/attribute"
)

type AttrBuff struct {
	Buff
	Attrs map[attribute.AttrType]int32
}

// NewAttrBuff creates a new attribute buff
func NewAttrBuff(modelId int32) *AttrBuff {
	buffData := config.QueryById[configdomain.BuffData](modelId)

	attr_map := make(map[attribute.AttrType]int32)
	// attack_10;defense_11
	attrStrs := strings.Split(buffData.Params, ";")
	for _, attrStr := range attrStrs {
		attrType := strings.Split(attrStr, "_")[0]
		attrValue := strings.Split(attrStr, "_")[1]
		attr_map[attribute.AttrType(attrType)], _ = util.StringToInt32(attrValue)
	}
	return &AttrBuff{
		Buff: Buff{
			ModelId:  modelId,
			Id:       util.GetNextID(),
			Duration: buffData.Duration,
		},
		Attrs: attr_map,
	}
}
