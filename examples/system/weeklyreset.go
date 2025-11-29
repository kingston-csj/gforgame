package system

import (
	"io/github/gforgame/examples/context"
	"strconv"
)

type WeeklyReset struct {
	ID        string `json:"id"`
	ResetTime int64  `json:"reset_time"`
	Data      interface{}
}

// DoParse 方法用于解析数据
func (d *WeeklyReset) DoParse() interface{} {		
	data := d.loadFromDb()
	if data == "" {
		return int64(0)
	}
	value, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return int64(0)
	}
	return value
}

// DoSave 方法用于保存数据
func (d *WeeklyReset) DoSave() string {
	return strconv.FormatInt(d.ResetTime, 10)
}

// GetID 方法用于获取参数 ID
func (d *WeeklyReset) GetID() string {
	return d.ID
}

// GetValue 方法用于获取参数值
func (d *WeeklyReset) GetValue() interface{} {
	if d.Data == nil {
		d.Data = d.DoParse()
	}
	return d.Data
}

// Save 方法用于保存参数
func (d *WeeklyReset) Save(data interface{}) {
	d.ResetTime = data.(int64)
	cache, _ := context.CacheManager.GetCache("systemparameter")
	cache.Set(d.GetID(), d)
	record := GetSystemParameterService().GetOrCreateSystemParameterRecord("1003")
	record.Data = d.DoSave()
	context.DbService.SaveToDb(record)
}

// loadFromDb 方法用于从数据库加载数据
func (d *WeeklyReset) loadFromDb() string {
	record := GetSystemParameterService().GetOrCreateSystemParameterRecord("1003")
	if record == nil {
		return ""
	}
	return record.GetData()
}
