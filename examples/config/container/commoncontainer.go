package container

import (
	"io/github/gforgame/data"
	configdomain "io/github/gforgame/examples/domain/config"
	"strconv"
)

type CommonContainer struct {
	*data.Container[int32, configdomain.CommonData]
	keys map[string]string
}

func (c *CommonContainer) Init() {
	if c.Container == nil {
		c.Container = data.NewContainer[int32, configdomain.CommonData]()
	}
	c.Container.Init()
}

// AfterLoad 数据加载后的处理
func (c *CommonContainer) AfterLoad() {
	c.keys = make(map[string]string)
	for _, record := range c.Container.GetAllRecords() {
		c.keys[record.Key] = record.Value
	}
}

// GetValue 获取配置值
func (c *CommonContainer) GetStringValue(key string) string {
	return c.keys[key]
}
   // 转换为int32
func (c *CommonContainer) GetInt32Value(key string) int32 {
	value, _ := strconv.ParseInt(c.keys[key], 10, 32)
	return int32(value)
}

func (c *CommonContainer) GetFloat32Value(key string) float32 {
	value, _ := strconv.ParseFloat(c.keys[key], 32)
	return float32(value)
}