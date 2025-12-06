package container

import (
	"io/github/gforgame/data"
	configdomain "io/github/gforgame/examples/domain/config"
)

type VipContainer struct {
	*data.Container[int32, configdomain.VipData]
}

func (c *VipContainer) Init() {
	if c.Container == nil {
		c.Container = data.NewContainer[int32, configdomain.VipData]()
	}
	c.Container.Init()
}

// AfterLoad 数据加载后的处理
func (c *VipContainer) AfterLoad() {

}

  