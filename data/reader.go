package data

type DataReader interface {
	Read(filePath string, clazz any) ([]any, error)
}
