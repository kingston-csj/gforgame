package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/conv"

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
		if rowLine < 3  {
			continue
		}
		if rowLine == 3 {
			headers, err = r.readHeader(clazz, row.Cells)
			if err != nil {
				return nil, err
			}
			continue
		}
		firstCell := getCellValue(row.Cells[0])
		if conv.EqualsIgnoreCase(firstCell, "") {
			break
		}

		// if len(headers) == 0 {
		// 	continue
		// }

		record := r.readExcelRow(headers, row)
		records = append(records, record)

		if conv.EqualsIgnoreCase(firstCell, "") {
			break
		}
	}

	return r.readRecords(clazz, records)
}

func resolveExcelFilePath(filePath string) string {
	if filepath.IsAbs(filePath) {
		return filePath
	}
	// 优先从 exe 目录向上查找，兼容打包和 IDE 的 __debug_bin。
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		if path, ok := findExcelFileFromBase(exeDir, filePath); ok {
			return path
		}
		if root, ok := findProjectRoot(exeDir); ok {
			return filepath.Join(root, "config", "excel", filePath)
		}
	}
	// 再从当前工作目录向上查找，兼容 `go test` 等场景。
	if cwd, err := os.Getwd(); err == nil {
		if path, ok := findExcelFileFromBase(cwd, filePath); ok {
			return path
		}
		if root, ok := findProjectRoot(cwd); ok {
			return filepath.Join(root, "config", "excel", filePath)
		}
	}
	// 兜底：保持原有相对路径行为。
	if abs, err := filepath.Abs(filepath.Join("config", "excel", filePath)); err == nil {
		return abs
	}
	return filepath.Join("config", "excel", filePath)
}

func findExcelFileFromBase(baseDir, filePath string) (string, bool) {
	dir := filepath.Clean(baseDir)
	for {
		candidate := filepath.Join(dir, "config", "excel", filePath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

func findProjectRoot(baseDir string) (string, bool) {
	dir := filepath.Clean(baseDir)
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

func (r *ExcelDataReader) readRecords(clazz any, rows [][]CellColumn) ([]any, error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("",fmt.Errorf("readRecords panic recovered: %v", err))
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
		if i > len(headers) {
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

// 处理分号分隔的切片数?
func parseSliceValue(value string, fieldType reflect.Type) (interface{}, error) {
	strValues := strings.Split(value, ";")
	sliceVal := reflect.MakeSlice(fieldType, len(strValues), len(strValues))

	// 根据切片的元素类型进行转换
		elemType := fieldType.Elem()
	for i, strVal := range strValues {
		strVal = strings.TrimSpace(strVal)

		var elemVal reflect.Value
		switch elemType.Kind() {
		case reflect.String:
			elemVal = reflect.ValueOf(strVal)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if strVal == "" {
				elemVal = reflect.Zero(elemType)
				break
			}
			num, err := strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse slice element '%s' as %s: %w", strVal, elemType.Kind(), err)
			}
			elemVal = reflect.ValueOf(num).Convert(elemType)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if strVal == "" {
				elemVal = reflect.Zero(elemType)
				break
			}
			num, err := strconv.ParseUint(strVal, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse slice element '%s' as %s: %w", strVal, elemType.Kind(), err)
			}
			elemVal = reflect.ValueOf(num).Convert(elemType)
		case reflect.Float32, reflect.Float64:
			if strVal == "" {
				elemVal = reflect.Zero(elemType)
				break
			}
			num, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse slice element '%s' as %s: %w", strVal, elemType.Kind(), err)
			}
			elemVal = reflect.ValueOf(num).Convert(elemType)
		case reflect.Bool:
			if strVal == "" {
				elemVal = reflect.Zero(elemType)
				break
			}
			b, err := strconv.ParseBool(strVal)
			if err != nil {
				return nil, fmt.Errorf("failed to parse slice element '%s' as bool: %w", strVal, err)
			}
			elemVal = reflect.ValueOf(b)
		default:
			return nil, fmt.Errorf("unsupported slice element type: %s", elemType.Kind())
		}

		sliceVal.Index(i).Set(elemVal)
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
		// 处理嵌套 JSON 对象
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
		tag := field.Tag.Get("excel") // 获取 Tag ?
		if strings.EqualFold(tag, tagValue) { // 忽略大小写匹?
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

type CellHeader struct {
	Column string
	Field  reflect.Value
}

type CellColumn struct {
	Header CellHeader
	Value  string
}
