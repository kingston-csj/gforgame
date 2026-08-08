package system

import systemrepo "github.com/forfun/gforgame/internal/infra/repository/system"

type WeeklyReset struct {
	baseInt64Parameter
	ResetTime int64 `json:"reset_time"`
	Data      any
}

func NewWeeklyReset(repo *systemrepo.SystemRepository) *WeeklyReset {
	w := &WeeklyReset{}
	w.baseInt64Parameter.init(SystemParamIDWeeklyReset, repo)
	return w
}

// DoParse 方法用于解析数据
func (d *WeeklyReset) DoParse() any {
	value := d.baseInt64Parameter.parseFromStore()
	d.ResetTime = value
	return value
}

// DoSave 方法用于保存数据
func (d *WeeklyReset) DoSave() string {
	return formatInt64(d.ResetTime)
}

// GetID 方法用于获取参数 ID
func (d *WeeklyReset) GetID() string {
	return d.baseInt64Parameter.getID()
}

// GetValue 方法用于获取参数值
func (d *WeeklyReset) GetValue() any {
	v := d.baseInt64Parameter.getValue()
	d.ResetTime = v
	return v
}

// Save 方法用于保存参数
func (d *WeeklyReset) Save(data any) {
	d.ResetTime = data.(int64)
	d.baseInt64Parameter.saveValue(d.ResetTime)
}
