package bootstrap

import (
	"goblog/pkg/model"
	"time"
)

// SetupDB 初始化数据库连接
func SetupDB() {

	db := model.ConnectDB()
	sqlDB, _ := db.DB()

	// 最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 最大空闲连接数
	sqlDB.SetMaxIdleConns(25)
	// 每个连接的过期时间
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

}
