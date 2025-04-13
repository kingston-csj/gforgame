package utils

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

func IsEmpty(s string) bool {
	return s == "" || s == "null"
}

func IsBlank(s string) bool {
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
