package system

type SystemParameter interface {
	DoParse() interface{}
	DoSave() string
	GetID() int
	GetValue() interface{}
	Save(data interface{})
}
