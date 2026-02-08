package container

import (
	"io/github/gforgame/data"
	configdomain "io/github/gforgame/examples/domain/config"
	"math/rand"
)

type NameContainer struct {
	*data.Container[int32, configdomain.NameData]
	allFirstNames []string
	allLastNames  []string
}

func (c *NameContainer) Init() {
	if c.Container == nil {
		c.Container = data.NewContainer[int32, configdomain.NameData]()
	}
	c.allFirstNames = make([]string, 0)
	c.allLastNames = make([]string, 0)
	c.Container.Init()
}

// AfterLoad 数据加载后的处理
func (c *NameContainer) AfterLoad() {
	for _, record := range c.GetAllRecords() {
		c.allFirstNames = append(c.allFirstNames, record.First)
		c.allLastNames = append(c.allLastNames, record.Last)
	}
}

func (c *NameContainer) GetRandomName() string {
	first := c.allFirstNames[rand.Intn(len(c.allFirstNames))]
	last := c.allLastNames[rand.Intn(len(c.allLastNames))]
	return first + last
}
