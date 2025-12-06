package util

import (
	"fmt"
	"strconv"
	"strings"
)

// -------------------------- Boolean 转换 --------------------------
// BooleanValue 将 any 类型转换为 bool（默认值：false）
func BooleanValue(object any) bool {
	return BooleanValueWithDefault(object, false)
}

// BooleanValueWithDefault 将 any 类型转换为 bool（支持自定义默认值）
func BooleanValueWithDefault(object any, defaultValue bool) bool {
	if object == nil {
		return defaultValue
	}

	// 直接匹配 bool 类型（Go 无包装类，直接断言基础类型）
	if b, ok := object.(bool); ok {
		return b
	}

	// 其他类型转字符串后解析（兼容字符串、数值等）
	str := toString(object)
	lowerStr := strings.ToLower(str)
	switch lowerStr {
	case "true", "1", "yes", "y":
		return true
	case "false", "0", "no", "n":
		return false
	default:
		return defaultValue // 解析失败返回默认值
	}
}

// -------------------------- Byte 转换（Go 中为 uint8） --------------------------
// ByteValue 将 any 类型转换为 byte（默认值：0）
func ByteValue(object any) uint8 {
	return ByteValueWithDefault(object, 0)
}

// ByteValueWithDefault 将 any 类型转换为 byte（支持自定义默认值）
func ByteValueWithDefault(object any, defaultValue uint8) uint8 {
	if object == nil {
		return defaultValue
	}

	// 直接匹配 uint8（Go 的 byte 等价于 uint8）
	if b, ok := object.(uint8); ok {
		return b
	}

	// 其他数值类型先转换为 int，再强转 byte（避免溢出）
	switch v := object.(type) {
	case int:
		return uint8(v)
	case int16:
		return uint8(v)
	case int32:
		return uint8(v)
	case int64:
		return uint8(v)
	case float32:
		return uint8(v)
	case float64:
		return uint8(v)
	case uint16:
		return uint8(v)
	case uint32:
		return uint8(v)
	case uint64:
		return uint8(v)
	}

	// 字符串类型解析
	str := toString(object)
	val, err := strconv.ParseUint(str, 10, 8) // 解析为 8 位无符号整数
	if err != nil {
		return defaultValue
	}
	return uint8(val)
}

// -------------------------- Short 转换（Go 中为 int16） --------------------------
// ShortValue 将 any 类型转换为 int16（默认值：0）
func ShortValue(object any) int16 {
	return ShortValueWithDefault(object, 0)
}

// ShortValueWithDefault 将 any 类型转换为 int16（支持自定义默认值）
func ShortValueWithDefault(object any, defaultValue int16) int16 {
	if object == nil {
		return defaultValue
	}

	// 直接匹配 int16（Go 的 short 等价于 int16）
	if s, ok := object.(int16); ok {
		return s
	}

	// 其他数值类型强转（避免溢出）
	switch v := object.(type) {
	case int:
		return int16(v)
	case uint8:
		return int16(v)
	case int32:
		return int16(v)
	case int64:
		return int16(v)
	case float32:
		return int16(v)
	case float64:
		return int16(v)
	case uint16:
		return int16(v)
	}

	// 字符串类型解析
	str := toString(object)
	val, err := strconv.ParseInt(str, 10, 16) // 解析为 16 位整数
	if err != nil {
		return defaultValue
	}
	return int16(val)
}

// -------------------------- Int 转换（Go 中为 int） --------------------------
// IntValue 将 any 类型转换为 int（默认值：0）
func IntValue(object any) int {
	return IntValueWithDefault(object, 0)
}

// IntValueWithDefault 将 any 类型转换为 int（支持自定义默认值）
func IntValueWithDefault(object any, defaultValue int) int {
	if object == nil {
		return defaultValue
	}

	// 直接匹配 int 类型
	if i, ok := object.(int); ok {
		return i
	}

	// 其他数值类型强转
	switch v := object.(type) {
	case uint8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	}

	// 字符串类型解析
	str := toString(object)
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return val
}

// -------------------------- Long 转换（Go 中为 int64） --------------------------
// LongValue 将 any 类型转换为 int64（默认值：0）
func LongValue(object any) int64 {
	return LongValueWithDefault(object, 0)
}

// LongValueWithDefault 将 any 类型转换为 int64（支持自定义默认值）
func LongValueWithDefault(object any, defaultValue int64) int64 {
	if object == nil {
		return defaultValue
	}

	// 直接匹配 int64 类型
	if l, ok := object.(int64); ok {
		return l
	}

	// 其他数值类型强转
	switch v := object.(type) {
	case uint8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	}

	// 字符串类型解析
	str := toString(object)
	val, err := strconv.ParseInt(str, 10, 64) // 解析为 64 位整数
	if err != nil {
		return defaultValue
	}
	return val
}

// -------------------------- Double 转换（Go 中为 float64） --------------------------
// DoubleValue 将 any 类型转换为 float64（默认值：0.0）
func DoubleValue(object any) float64 {
	return DoubleValueWithDefault(object, 0.0)
}

// DoubleValueWithDefault 将 any 类型转换为 float64（支持自定义默认值）
func DoubleValueWithDefault(object any, defaultValue float64) float64 {
	if object == nil {
		return defaultValue
	}

	// 直接匹配 float64 类型（Go 的 double 等价于 float64）
	if d, ok := object.(float64); ok {
		return d
	}

	// 其他数值类型强转
	switch v := object.(type) {
	case uint8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case float32:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	}

	// 字符串类型解析
	str := toString(object)
	val, err := strconv.ParseFloat(str, 64) // 解析为 64 位浮点数
	if err != nil {
		return defaultValue
	}
	return val
}

// -------------------------- 内部辅助函数 --------------------------
// toString 将 any 类型转换为字符串（兼容各种类型）
func toString(object any) string {
	if object == nil {
		return ""
	}

	// 直接匹配 string 类型
	if s, ok := object.(string); ok {
		return s
	}

	// 其他类型通过 fmt 格式化转为字符串（兼容数值、布尔等）
	return fmt.Sprintf("%v", object)
}