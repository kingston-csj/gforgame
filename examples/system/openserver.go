package system

import (
	"io/github/gforgame/examples/context"
)

type OpenSeverTime struct {
	ID        string `json:"id"`
	Date string  `json:"date"`
	Data interface{}
}

// DoParse 方法用于解析数据
func (d *OpenSeverTime) DoParse() interface{} {
	return d.loadFromDb()
}

// DoSave 方法用于保存数据
func (d *OpenSeverTime) DoSave() string {
	return d.Date
}

// GetID 方法用于获取参数 ID
func (d *OpenSeverTime) GetID() string {
	return d.ID
}

// GetValue 方法用于获取参数值
func (d *OpenSeverTime) GetValue() interface{} {
	if d.Data == nil {
		d.Data = d.DoParse()
	}
	return d.Data
}

// Save 方法用于保存参数
func (d *OpenSeverTime) Save(data interface{}) {
	d.Date = data.(string)
	cache, _ := context.CacheManager.GetCache("systemparameter")
	cache.Set(d.GetID(), d)
	record := GetSystemParameterService().GetOrCreateSystemParameterRecord("1004")
	record.Data = d.DoSave()
	context.DbService.SaveToDb(record)
}

// loadFromDb 方法用于从数据库加载数据
func (d *OpenSeverTime) loadFromDb() string {
	record := GetSystemParameterService().GetOrCreateSystemParameterRecord("1004")
	if record == nil {
		return ""
	}
	return record.GetData()
}
