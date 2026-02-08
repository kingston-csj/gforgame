package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// -------------------- 核心配置（可选参数） --------------------
// JSON序列化选项（用于自定义规则，避免函数参数过多）
type MarshalOptions struct {
	Indent      bool   // 是否格式化缩进（便于阅读）
	IgnoreEmpty bool   // 是否忽略空值字段（需配合 struct tag `omitempty` 使用）
	TimeFormat  string // 时间类型（time.Time）的序列化格式，默认RFC3339
}

// 默认选项（无缩进、不忽略空值、时间用RFC3339格式）
var defaultOptions = MarshalOptions{
	Indent:      false,
	IgnoreEmpty: false,
	TimeFormat:  time.RFC3339,
}

// -------------------- 核心函数 --------------------

// StructToJSON 结构体转JSON字符串（使用默认配置）
// 入参：任意结构体（或指针）
// 返回：JSON字符串 / 错误（如循环引用、非导出字段等）
func StructToJSON(v any) (string, error) {
	return StructToJSONWithOptions(v, defaultOptions)
}

// StructToPrettyJSON 结构体转格式化的JSON字符串（带缩进，便于阅读）
func StructToPrettyJSON(v any) (string, error) {
	options := defaultOptions
	options.Indent = true
	return StructToJSONWithOptions(v, options)
}

// StructToJSONWithOptions 自定义选项的结构体转JSON（灵活配置）
func StructToJSONWithOptions(v any, options MarshalOptions) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	
	// 禁用HTML转义（避免<、>、&被转义为\u003c等）
	// 这是使用 json.Encoder 而非 json.Marshal 的主要原因
	encoder.SetEscapeHTML(false)

	// 处理缩进
	if options.Indent {
		encoder.SetIndent("", "  ") // 缩进2个空格
	}

	if err := encoder.Encode(v); err != nil {
		return "", fmt.Errorf("struct转JSON失败: %w", err)
	}

	// json.Encoder.Encode 会自动在末尾添加换行符，为了保持与 json.Marshal 一致的行为，我们需要去掉它
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}

	return string(result), nil
}

// JsonToStruct JSON字符串转结构体
// 入参：jsonStr (JSON字符串), v (接收结果的结构体指针)
func JsonToStruct(jsonStr string, v any) error {
	if jsonStr == "" {
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), v)
}

// JsonBytesToStruct JSON字节切片转结构体
// 入参：data (JSON字节切片), v (接收结果的结构体指针)
func JsonBytesToStruct(data []byte, v any) error {
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}

// -------------------- 辅助函数（常用场景） --------------------

// MustStructToJSON 忽略错误的结构体转JSON（仅用于确定无错误的场景，如测试）
func MustStructToJSON(v any) string {
	str, err := StructToJSON(v)
	if err != nil {
		panic(fmt.Sprintf("StructToJSON panic: %v", err))
	}
	return str
}
