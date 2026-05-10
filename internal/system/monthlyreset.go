package system

type MonthlyReset struct {
	baseInt64Parameter
	ResetTime int64 `json:"reset_time"`
	Data      interface{}
}

func NewMonthlyReset() *MonthlyReset {
	m := &MonthlyReset{}
	m.baseInt64Parameter.init(SystemParamIDMonthlyReset)
	return m
}

// DoParse 方法用于解析数据
func (d *MonthlyReset) DoParse() interface{} {
	value := d.baseInt64Parameter.parseFromStore(func() string {
		return d.loadFromDb()
	})
	d.ResetTime = value
	return value
}

// DoSave 方法用于保存数据
func (d *MonthlyReset) DoSave() string {
	return formatInt64(d.ResetTime)
}

// GetID 方法用于获取参数 ID
func (d *MonthlyReset) GetID() string {
	d.baseInt64Parameter.init(SystemParamIDMonthlyReset)
	return d.baseInt64Parameter.getID()
}

// GetValue 方法用于获取参数值
func (d *MonthlyReset) GetValue() interface{} {
	v := d.baseInt64Parameter.getValue(func() string {
		return d.loadFromDb()
	})
	d.ResetTime = v
	return v
}

// Save 方法用于保存参数
func (d *MonthlyReset) Save(data interface{}) {
	d.ResetTime = data.(int64)
	d.baseInt64Parameter.saveValue(d.ResetTime, d)
}

// loadFromDb 方法用于从数据库加载数据
func (d *MonthlyReset) loadFromDb() string {
	return loadSystemParameterValue(d.GetID())
}
