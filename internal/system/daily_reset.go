package system

import systemrepo "github.com/forfun/gforgame/internal/infra/repository/system"

type DailyReset struct {
	baseInt64Parameter
	ResetTime int64 `json:"reset_time"`
}

func NewDailyReset(repo *systemrepo.SystemRepository) *DailyReset {
	d := &DailyReset{}
	d.baseInt64Parameter.init(SystemParamIDDailyReset, repo)
	d.DoParse()
	return d
}

// DoParse 方法用于解析数据
func (d *DailyReset) DoParse() any {
	value := d.baseInt64Parameter.parseFromStore()
	d.ResetTime = value
	return value
}

// DoSave 方法用于保存数据
func (d *DailyReset) DoSave() string {
	return formatInt64(d.ResetTime)
}

// GetID 方法用于获取参数 ID
func (d *DailyReset) GetID() string {
	return d.baseInt64Parameter.getID()
}

// GetValue 方法用于获取参数值
func (d *DailyReset) GetValue() any {
	v := d.baseInt64Parameter.getValue()
	d.ResetTime = v
	return v
}

// Save 方法用于保存参数
func (d *DailyReset) Save(data any) {
	d.ResetTime = data.(int64)
	d.baseInt64Parameter.saveValue(d.ResetTime)
}
