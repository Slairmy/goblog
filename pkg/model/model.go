package model

import (
	"goblog/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	gormlogger "gorm.io/gorm/logger"
)

// DB 全局变量
var DB *gorm.DB

// ConnectDB 连接gormDB
func ConnectDB() *gorm.DB {
	var err error
	config := mysql.New(mysql.Config{
		DSN: "root:+Darlingmy521@tcp(127.0.0.1:3306)/goblog?charset=utf8&parseTime=True&loc=Local",
	})

	DB, err = gorm.Open(config, &gorm.Config{
		// gorm 打印日志
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	logger.LogError(err)

	return DB
}
