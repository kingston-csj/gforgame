package buff

import (
	"io/github/gforgame/examples/config"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/fight/attribute"
	"io/github/gforgame/util"
	"strings"
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
