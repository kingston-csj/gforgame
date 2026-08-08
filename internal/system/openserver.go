package system

import systemrepo "github.com/forfun/gforgame/internal/infra/repository/system"

type OpenSeverTime struct {
	baseStringParameter
	Date string `json:"date"`
	Data any
}

func NewOpenServerTime(repo *systemrepo.SystemRepository) *OpenSeverTime {
	o := &OpenSeverTime{}
	o.baseStringParameter.init(SystemParamIDOpenServer, repo)
	return o
}

// DoParse 方法用于解析数据
func (d *OpenSeverTime) DoParse() any {
	value := d.baseStringParameter.parseFromStore()
	d.Date = value
	return value
}

// DoSave 方法用于保存数据
func (d *OpenSeverTime) DoSave() string {
	return d.Date
}

// GetID 方法用于获取参数 ID
func (d *OpenSeverTime) GetID() string {
	return d.baseStringParameter.getID()
}

// GetValue 方法用于获取参数值
func (d *OpenSeverTime) GetValue() any {
	v := d.baseStringParameter.getValue()
	return v
}

// Save 方法用于保存参数
func (d *OpenSeverTime) Save(data any) {
	d.Date = data.(string)
	d.baseStringParameter.saveValue(d.DoSave())
}
