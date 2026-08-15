package data

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	util "github.com/forfun/gforgame/common/util/conv"
	"github.com/forfun/gforgame/common/util/pathutil"

	"github.com/tealeg/xlsx"
)

type ExcelDataReader struct {
	ignoreUnknownFields bool
}

func NewExcelDataReader(ignoreUnknownFields bool) *ExcelDataReader {
	return &ExcelDataReader{
		ignoreUnknownFields: ignoreUnknownFields,
	}
}

func (r *ExcelDataReader) Read(filePath string, clazz any) ([]any, error) {
	excelFilePath := resolveExcelFilePath(filePath)
	xlFile, err := xlsx.OpenFile(excelFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %v", err)
	}

	sheet := xlFile.Sheets[0]
	rows := sheet.Rows

	var headers []CellHeader
	var records [][]CellColumn
	rowLine := 0
	// 遍历每一行
	for _, row := range rows {
		rowLine++
		if rowLine < 3 {
			continue
		}
		if rowLine == 3 {
			headers, err = r.readHeader(row.Cells)
			if err != nil {
				return nil, err
			}
			continue
		}
		firstCell := getCellValue(row.Cells[0])
		if util.EqualsIgnoreCase(firstCell, "") {
			break
		}

		// if len(headers) == 0 {
		// 	continue
		// }

		record := r.readExcelRow(headers, row)
		records = append(records, record)

		if util.EqualsIgnoreCase(firstCell, "") {
			break
		}
	}

	return r.readRecords(clazz, records)
}

func resolveExcelFilePath(filePath string) string {
	if filepath.IsAbs(filePath) {
		return filePath
	}
	return pathutil.ResolveFilePath(filepath.Join("config", "excel", filePath))
}

func (r *ExcelDataReader) readRecords(clazz any, rows [][]CellColumn) (records []any, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("read records panic recovered: %v", recovered)
			records = nil
		}
	}()

	clazzType := reflect.TypeOf(clazz)
	if clazzType == nil || clazzType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("clazz must be a non-nil pointer, got %T", clazz)
	}
	clazzType = clazzType.Elem()

	for i, row := range rows {
		obj := reflect.New(clazzType).Elem()

		for _, column := range row {
			colName := column.Header.Column
			if colName == "" {
				continue
			}

			// 根据 Tag 查找字段
			field, err := findFieldByTag(obj, colName)
			if err != nil {
				if !r.ignoreUnknownFields {
					return nil, fmt.Errorf("row %d, column '%s': %v", i+1, colName, err)
				}
				continue
			}

			fieldVal, err := convertValue(column.Value, field.Type())
			if err != nil {
				return nil, fmt.Errorf("row %d, column '%s': %v", i+1, colName, err)
			}

			if !field.CanSet() {
				return nil, fmt.Errorf("row %d, column '%s': field cannot be set", i+1, colName)
			}
			convertedValue := reflect.ValueOf(fieldVal)
			if !convertedValue.Type().AssignableTo(field.Type()) {
				if convertedValue.Type().ConvertibleTo(field.Type()) {
					convertedValue = convertedValue.Convert(field.Type())
				} else {
					return nil, fmt.Errorf("row %d, column '%s': value type %v cannot be assigned to %v", i+1, colName, convertedValue.Type(), field.Type())
				}
			}
			field.Set(convertedValue)
		}

		records = append(records, obj.Interface())
	}

	return records, nil
}

func (r *ExcelDataReader) readHeader(cells []*xlsx.Cell) ([]CellHeader, error) {
	var headers []CellHeader

	for _, cell := range cells {
		cellValue := getCellValue(cell)
		// 空值表示列头结束
		if util.EqualsIgnoreCase(cellValue, "") {
			break
		}
		header := CellHeader{
			Column: cellValue,
		}

		headers = append(headers, header)
	}

	return headers, nil
}

func getCellValue(cell *xlsx.Cell) string {
	if cell == nil {
		return ""
	}
	return cell.String()
}

func (r *ExcelDataReader) readExcelRow(headers []CellHeader, row *xlsx.Row) []CellColumn {
	var columns []CellColumn

	for i, cell := range row.Cells {
		if i >= len(headers) {
			break
		}

		cellValue := getCellValue(cell)
		column := CellColumn{
			Header: headers[i],
			Value:  cellValue,
		}
		columns = append(columns, column)
	}

	return columns
}

func parseScalarValue(value string, fieldType reflect.Type) (reflect.Value, error) {
	if value == "" {
		return reflect.Zero(fieldType), nil
	}

	switch fieldType.Kind() {
	case reflect.String:
		return reflect.ValueOf(value).Convert(fieldType), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.ParseInt(value, 10, fieldType.Bits())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("failed to parse int: value='%s', type=%v, error=%v", value, fieldType.Kind(), err)
		}
		val := reflect.New(fieldType).Elem()
		val.SetInt(num)
		return val, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := strconv.ParseUint(value, 10, fieldType.Bits())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("failed to parse uint: value='%s', type=%v, error=%v", value, fieldType.Kind(), err)
		}
		val := reflect.New(fieldType).Elem()
		val.SetUint(num)
		return val, nil
	case reflect.Float32, reflect.Float64:
		num, err := strconv.ParseFloat(value, fieldType.Bits())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("failed to parse float: value='%s', type=%v, error=%v", value, fieldType.Kind(), err)
		}
		val := reflect.New(fieldType).Elem()
		val.SetFloat(num)
		return val, nil
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("failed to parse bool: value='%s', error=%v", value, err)
		}
		val := reflect.New(fieldType).Elem()
		val.SetBool(b)
		return val, nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported scalar type: %v", fieldType.Kind())
	}
}

// 处理分号分隔的切片数值
func parseSliceValue(value string, fieldType reflect.Type) (any, error) {
	strValues := strings.Split(value, ";")
	sliceVal := reflect.MakeSlice(fieldType, len(strValues), len(strValues))

	// 根据切片的元素类型进行转换
	elemType := fieldType.Elem()
	for i, strVal := range strValues {
		strVal = strings.TrimSpace(strVal)
		elemVal, err := parseScalarValue(strVal, elemType)
		if err != nil {
			return nil, fmt.Errorf("failed to parse slice element '%s' as %s: %w", strVal, elemType.Kind(), err)
		}
		sliceVal.Index(i).Set(elemVal)
	}
	return sliceVal.Interface(), nil
}

func parseArrayValue(value string, fieldType reflect.Type) (any, error) {
	strValues := strings.Split(value, ";")
	arrayVal := reflect.New(fieldType).Elem()
	elemType := fieldType.Elem()
	if len(strValues) > fieldType.Len() {
		return nil, fmt.Errorf("array element count %d exceeds fixed length %d", len(strValues), fieldType.Len())
	}
	for i, strVal := range strValues {
		elemVal, err := parseScalarValue(strings.TrimSpace(strVal), elemType)
		if err != nil {
			return nil, fmt.Errorf("failed to parse array element '%s' as %s: %w", strVal, elemType.Kind(), err)
		}
		arrayVal.Index(i).Set(elemVal)
	}
	return arrayVal.Interface(), nil
}

func convertValue(value string, fieldType reflect.Type) (any, error) {
	switch fieldType.Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		val, err := parseScalarValue(value, fieldType)
		if err != nil {
			return nil, err
		}
		return val.Interface(), nil
	case reflect.Array:
		return parseArrayValue(value, fieldType)
	case reflect.Slice, reflect.Struct:
		fieldVal := reflect.New(fieldType).Interface()
		if err := json.Unmarshal([]byte(value), &fieldVal); err != nil {
			if fieldType.Kind() == reflect.Slice {
				return parseSliceValue(value, fieldType)
			}
			return nil, fmt.Errorf("failed to unmarshal JSON and not a simple slice: %v", err)
		}
		return reflect.ValueOf(fieldVal).Elem().Interface(), nil
	case reflect.Interface:
		if result, err := ConvertRegistryInstance.Convert(value, fieldType); err != nil {
			return nil, err
		} else if result != nil {
			return result, nil
		}
		return nil, fmt.Errorf("no converter registered for interface type: %v", fieldType)
	default:
		return nil, fmt.Errorf("unsupported type: %v", fieldType.Kind())
	}
}

// 根据 Tag 查找字段
func findFieldByTag(obj reflect.Value, tagValue string) (reflect.Value, error) {
	// 如果传入的是指针，解引用
	if obj.Kind() == reflect.Ptr {
		if obj.IsNil() {
			return reflect.Value{}, fmt.Errorf("nil pointer")
		}
		obj = obj.Elem()
	}

	objType := obj.Type()
	if objType.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("not a struct")
	}

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		tag := field.Tag.Get("excel")         // 获取 Tag 值
		if strings.EqualFold(tag, tagValue) { // 忽略大小写匹配
			return obj.Field(i), nil
		}

		// 递归查找匿名结构体（嵌入字段）
		if field.Anonymous {
			fieldVal := obj.Field(i)

			// 如果是指针且nil，需要初始化
			if fieldVal.Kind() == reflect.Ptr {
				if fieldVal.IsNil() {
					if fieldVal.CanSet() {
						newValue := reflect.New(fieldVal.Type().Elem())
						fieldVal.Set(newValue)
					}
				}
			}

			// 递归调用
			if res, err := findFieldByTag(fieldVal, tagValue); err == nil {
				return res, nil
			}
		}
	}
	return reflect.Value{}, fmt.Errorf("field with tag %s not found", tagValue)
}
