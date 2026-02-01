package util

import (
	"fmt"
	"strconv"
	"strings"
)

// ToIntIntMap
// 功能：将分隔符分割的字符串解析为有序的 map[int]int
// 参数：
//   params - 待解析的字符串（如 "1:10,2:20,3:30"）
//   unitDelimiter - 键值对之间的分隔符（如 ","）
//   valueDelimiter - 键和值之间的分隔符（如 ":"）
// 返回：
//  键值对 map（模拟 LinkedHashMap）、错误信息
func ToIntIntMap(params, unitDelimiter, valueDelimiter string) (map[int32]int32, error) {
	// 空字符串直接返回空map
	if params == "" {
		return make(map[int32]int32), nil
	}

	// 1. 用 unitDelimiter 分割为键值对片段
	splits := strings.Split(params, unitDelimiter)
	// 有序map：用普通map存键值，同时用切片存插入顺序
	result := make(map[int32]int32, len(splits))

	// 2. 遍历每个键值对片段
	for idx, split := range splits {
		// 跳过空片段（比如 params 末尾有分隔符的情况）
		if split == "" {
			continue
		}

		// 3. 用 valueDelimiter 分割键和值
		unit := strings.Split(split, valueDelimiter)
		// 校验分割后长度（必须是2，否则解析失败）
		if len(unit) != 2 {
			return nil, fmt.Errorf("第 %d 个片段解析失败：%s（分割后长度=%d，期望=2）", idx+1, split, len(unit))
		}

		// 4. 字符串转 int
		keyStr, valStr := unit[0], unit[1]
		key, err := strconv.Atoi(keyStr)
		if err != nil {
			return nil, fmt.Errorf("第 %d 个片段的键转int失败：%s，错误：%w", idx+1, keyStr, err)
		}
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return nil, fmt.Errorf("第 %d 个片段的值转int失败：%s，错误：%w", idx+1, valStr, err)
		}

		// 5. 放入map
		result[int32(key)] = int32(val)
	}

	return result, nil
}