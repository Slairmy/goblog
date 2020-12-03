package types

import (
	"goblog/pkg/logger"
	"strconv"
)

// Int64ToString 格式转换--int64 to string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// StringToInt 格式转换--string to int64
func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		logger.LogError(err)
	}
	return i
}

// Uint64ToString uint64 to string
func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}
