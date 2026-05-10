package system

import "strconv"

func formatInt64(v int64) string {
	return strconv.FormatInt(v, 10)
}
