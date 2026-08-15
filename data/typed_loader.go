package data

import "fmt"

// ReadTyped 以强类型方式读取 Excel 记录，返回 []*T。
func ReadTyped[T any](reader DataReader, filePath string) ([]*T, error) {
	raw, err := reader.Read(filePath, new(T))
	if err != nil {
		return nil, err
	}

	records := make([]*T, 0, len(raw))
	for i, item := range raw {
		switch v := item.(type) {
		case T:
			val := v
			records = append(records, &val)
		case *T:
			records = append(records, v)
		default:
			return nil, fmt.Errorf("row %d has incompatible type %T", i+1, item)
		}
	}
	return records, nil
}

// ProcessTableTyped 以强类型方式构建配置容器，避免调用方使用 any/反射。
func ProcessTableTyped[T any](
	reader DataReader,
	tableName string,
	filePath string,
	getID func(*T) int32,
	indexFuncs map[string]func(*T) any,
) (*Container[int32, T], error) {
	records, err := ReadTyped[T](reader, filePath)
	if err != nil {
		return nil, err
	}
	return BuildContainer(tableName, records, getID, indexFuncs)
}
