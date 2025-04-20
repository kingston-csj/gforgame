package util

import (
	"strconv"
	"strings"
)

func StringToInt32(s string) (int32, error) {
	num, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(num), nil
}

func StringToFloat32(s string) (float32, error) {
	num, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(num), nil
}

func IsEmptyString(s string) bool {
	return s == "" || s == "null"
}

func IsBlankString(s string) bool {
	return s == "" || s == "null" || strings.TrimSpace(s) == ""
}

func EqualsIgnoreCase(s1, s2 string) bool {
	if s1 == "" && s2 == "" {
		return true
	}
	if s1 == "" || s2 == "" {
		return false
	}
	return strings.EqualFold(s1, s2)
}
