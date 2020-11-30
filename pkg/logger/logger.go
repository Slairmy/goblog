package logger

import "log"

// LogError 打印错误日志
func LogError(err error) {
	if err != nil {
		log.Println(err)
	}

	// 打印日志到K8s
}
