package utils

import (
	"strconv"
)

func StringToInt32(s string) (int32, error) {
	num, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(num), nil
}
