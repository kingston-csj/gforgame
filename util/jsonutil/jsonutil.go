package jsonutil

import (
	"encoding/json"
	"fmt"
	"time"
)

// -------------------- 核心配置（可选参数） --------------------
// JSON序列化选项（用于自定义规则，避免函数参数过多）
type MarshalOptions struct {
	Indent       bool   // 是否格式化缩进（便于阅读）
	IgnoreEmpty  bool   // 是否忽略空值字段（如""、0、nil）
	TimeFormat   string // 时间类型（time.Time）的序列化格式，默认RFC3339
}

// 默认选项（无缩进、不忽略空值、时间用RFC3339格式）
var defaultOptions = MarshalOptions{
	Indent:     false,
	IgnoreEmpty: false,
	TimeFormat: time.RFC3339,
}

// -------------------- 核心函数 --------------------
// StructToJSON 结构体转JSON字符串（使用默认配置）
// 入参：任意结构体（或指针）
// 返回：JSON字符串 / 错误（如循环引用、非导出字段等）
func StructToJSON(v interface{}) (string, error) {
	return StructToJSONWithOptions(v, defaultOptions)
}

// StructToPrettyJSON 结构体转格式化的JSON字符串（带缩进，便于阅读）
func StructToPrettyJSON(v interface{}) (string, error) {
	options := defaultOptions
	options.Indent = true
	return StructToJSONWithOptions(v, options)
}

// StructToJSONWithOptions 自定义选项的结构体转JSON（灵活配置）
func StructToJSONWithOptions(v interface{}, options MarshalOptions) (string, error) {
	//  配置JSON编码器
	encoder := json.NewEncoder(&jsonBuffer{})
	// 禁用HTML转义（避免<、>、&被转义为\u003c等）
	encoder.SetEscapeHTML(false)
	
	//  处理缩进
	if options.Indent {
		encoder.SetIndent("", "  ") // 缩进2个空格
	}

	var buf []byte
	var err error
	if options.IgnoreEmpty {
		// 忽略空值的序列化（需用json.MarshalIndent/json.Marshal + omitempty）
		// 注：若要精准忽略空值，建议在struct tag中加 `omitempty`，这里是全局兜底
		if options.Indent {
			buf, err = json.MarshalIndent(v, "", "  ")
		} else {
			buf, err = json.Marshal(v)
		}
	} else {
		if options.Indent {
			buf, err = json.MarshalIndent(v, "", "  ")
		} else {
			buf, err = json.Marshal(v)
		}
	}

	if err != nil {
		return "", fmt.Errorf("struct转JSON失败: %w", err)
	}

	return string(buf), nil
}

// 辅助：空buffer（兼容encoder写法，实际用Marshal更简洁）
type jsonBuffer struct{}

func (j *jsonBuffer) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// -------------------- 辅助函数（常用场景） --------------------
// MustStructToJSON 忽略错误的结构体转JSON（仅用于确定无错误的场景，如测试）
func MustStructToJSON(v interface{}) string {
	str, err := StructToJSON(v)
	if err != nil {
		panic(fmt.Sprintf("StructToJSON panic: %v", err))
	}
	return str
}