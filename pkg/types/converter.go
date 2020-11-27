package types

import "strconv"

// Int64ToString 格式转换--int64 to string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}
