package data

import (
	"encoding/json"
	"fmt"
	"io/github/gforgame/util"
	"reflect"
	"strconv"
	"strings"

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
	// 使用 xlsx.OpenFile 打开 Excel 文件
	xlFile, err := xlsx.OpenFile("config/excel/" + filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %v", err)
	}

	sheet := xlFile.Sheets[0]
	rows := sheet.Rows

	var headers []CellHeader
	var records [][]CellColumn

	// 遍历每一行
	for _, row := range rows {
		firstCell := getCellValue(row.Cells[0])
		if util.EqualsIgnoreCase(firstCell, "HEADER") {
			headers, err = r.readHeader(clazz, row.Cells)
			if err != nil {
				return nil, err
			}
			continue
		}

		if len(headers) == 0 {
			continue
		}

		record := r.readExcelRow(headers, row)
		records = append(records, record)

		if util.EqualsIgnoreCase(firstCell, "END") {
			break
		}
	}

	return r.readRecords(clazz, records)
}

func (r *ExcelDataReader) readRecords(clazz any, rows [][]CellColumn) ([]any, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic:", clazz, rows)
		}
	}()

	var records []any
	clazzType := reflect.TypeOf(clazz).Elem()

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

			field.Set(reflect.ValueOf(fieldVal))
		}

		records = append(records, obj.Interface())
	}

	return records, nil
}

func (r *ExcelDataReader) readHeader(clazz any, cells []*xlsx.Cell) ([]CellHeader, error) {
	var headers []CellHeader

	for _, cell := range cells {
		cellValue := getCellValue(cell)
		if util.EqualsIgnoreCase(cellValue, "HEADER") {
			continue
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
		// 忽略 header 所在的第一列
		if i == 0 {
			continue
		}
		if i > len(headers) {
			break
		}

		cellValue := getCellValue(cell)
		column := CellColumn{
			// headers 从 0 开始，所以这里 -1
			Header: headers[i-1],
			Value:  cellValue,
		}
		columns = append(columns, column)
	}

	return columns
}

// 处理分号分隔的切片数据
func parseSliceValue(value string, fieldType reflect.Type) (interface{}, error) {
	strValues := strings.Split(value, ";")
	sliceVal := reflect.MakeSlice(fieldType, len(strValues), len(strValues))

	// 根据切片的元素类型进行转换
	elemType := fieldType.Elem()
	for i, strVal := range strValues {
		var elemVal interface{}
		var err error

		switch elemType.Kind() {
		case reflect.Int32:
			if num, err := strconv.ParseInt(strVal, 10, 32); err == nil {
				elemVal = int32(num)
			}
		case reflect.Int64:
			elemVal, err = strconv.ParseInt(strVal, 10, 64)
		case reflect.Float32:
			if num, err := strconv.ParseFloat(strVal, 32); err == nil {
				elemVal = float32(num)
			}
		case reflect.Float64:
			elemVal, err = strconv.ParseFloat(strVal, 64)
		default:
			elemVal = strVal
		}

		if err != nil {
			return nil, fmt.Errorf("failed to parse slice element: %v", err)
		}

		sliceVal.Index(i).Set(reflect.ValueOf(elemVal))
	}
	return sliceVal.Interface(), nil
}

func convertValue(value string, fieldType reflect.Type) (any, error) {
	switch fieldType.Kind() {
	case reflect.String:
		return value, nil
	case reflect.Int8, reflect.Int16, reflect.Int32:
		if value == "" {
			return int32(0), nil
		}
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse int: value='%s', type=%v, error=%v", value, fieldType.Kind(), err)
		}
		return int32(num), nil
	case reflect.Int:
		if value == "" {
			return int32(0), nil
		}
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse int: value='%s', type=%v, error=%v", value, fieldType.Kind(), err)
		}
		return int32(num), nil
	case reflect.Float32, reflect.Float64:
		if value == "" {
			return float32(0), nil
		}
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse float: value='%s', type=%v, error=%v", value, fieldType.Kind(), err)
		}
		return float32(num), nil
	case reflect.Bool:
		return strconv.ParseBool(value)
	case reflect.Array:
		return strings.Split(value, ";"), nil
	case reflect.Slice, reflect.Struct:
		// 处理嵌套的 JSON 对象
		fieldVal := reflect.New(fieldType).Interface()
		if err := json.Unmarshal([]byte(value), &fieldVal); err != nil {
			// 如果解析失败，尝试解析为数组
			if fieldType.Kind() == reflect.Slice {
				return parseSliceValue(value, fieldType)
			}
			return nil, fmt.Errorf("failed to unmarshal JSON and not a simple slice: %v", err)
		}
		return reflect.ValueOf(fieldVal).Elem().Interface(), nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", fieldType.Kind())
	}
}

// 根据 Tag 查找字段
func findFieldByTag(obj reflect.Value, tagValue string) (reflect.Value, error) {
	objType := obj.Type()
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		tag := field.Tag.Get("excel")         // 获取 Tag 值
		if strings.EqualFold(tag, tagValue) { // 忽略大小写匹配
			return obj.Field(i), nil
		}
	}
	return reflect.Value{}, fmt.Errorf("field with tag %s not found", tagValue)
}

type CellHeader struct {
	Column string
	Field  reflect.Value
}

type CellColumn struct {
	Header CellHeader
	Value  string
}
