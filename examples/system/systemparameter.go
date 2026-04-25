package system

type SystemParameter interface {
	DoParse() interface{}
	DoSave() string
	GetID() string
	GetValue() interface{}
	Save(data interface{})
}
