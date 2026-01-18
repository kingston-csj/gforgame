package container

import (
	"io/github/gforgame/data"
	"io/github/gforgame/examples/domain/config"
	"io/github/gforgame/util"
)
type GachaContainer struct {
	*data.Container[int32, config.GachaData]
}

func (c *GachaContainer) Init() {
	if c.Container == nil {
		c.Container = data.NewContainer[int32, config.GachaData]()
	}
	c.Container.Init()
}

func (c *GachaContainer) RandItem(gachaType int32) *config.GachaData {
	var pool []*config.GachaData
	var weights []int
	for _, record := range c.Container.Values() {
		if record.Type == gachaType {
			pool = append(pool, record)
			weights = append(weights, int(record.Weight))
		}
	}
	index,_ := util.RandomIndex(weights)
	return pool[index]
}