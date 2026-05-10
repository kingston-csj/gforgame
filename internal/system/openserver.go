package system

import (
	"sync"
	"sync/atomic"
)

type OpenSeverTime struct {
	baseStringParameter
	Date     string `json:"date"`
	Data     interface{}
	value    atomic.Value `json:"-"`
	loadOnce sync.Once    `json:"-"`
}

func NewOpenServerTime() *OpenSeverTime {
	o := &OpenSeverTime{}
	o.baseStringParameter.init(SystemParamIDOpenServer)
	return o
}

// DoParse 方法用于解析数据
func (d *OpenSeverTime) DoParse() interface{} {
	value := d.baseStringParameter.parseFromStore(func() string {
		return d.loadFromDb()
	})
	d.Date = value
	return value
}

// DoSave 方法用于保存数据
func (d *OpenSeverTime) DoSave() string {
	return d.Date
}

// GetID 方法用于获取参数 ID
func (d *OpenSeverTime) GetID() string {
	d.baseStringParameter.init(SystemParamIDOpenServer)
	return d.baseStringParameter.getID()
}

// GetValue 方法用于获取参数值
func (d *OpenSeverTime) GetValue() interface{} {
	v := d.baseStringParameter.getValue(func() string {
		return d.loadFromDb()
	})
	return v
}

// Save 方法用于保存参数
func (d *OpenSeverTime) Save(data interface{}) {
	d.Date = data.(string)
	d.baseStringParameter.saveValue(d.DoSave(), d)
}

// loadFromDb 方法用于从数据库加载数据
func (d *OpenSeverTime) loadFromDb() string {
	return loadSystemParameterValue(d.GetID())	
}
