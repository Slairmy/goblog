package config

import "goblog/pkg/config"

func init() {
	config.Add("app", config.StrMap{
		// 应用名称变量
		"name": config.Env("APP_NAME", "GoBlog"),

		// 环境信息变量
		"env": config.Env("APP_ENV", "production"),

		// 调试模式
		"debug": config.Env("APP_DEBUG", false),

		// 端口变量
		"port": config.Env("APP_PORT", "3000"),

		// cookie加密变量
		"key": config.Env("APP_KEY", "33446a9dcf9ea060a0a6532b166da32f304af0de"),
	})
}
