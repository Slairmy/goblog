package middlewares

import "net/http"

// HTTPHandlerFunc 简写 —— func(http.ResponseWriter, *http.Request)
// 定义作用与单个路由的中间件
type HTTPHandlerFunc func(http.ResponseWriter, *http.Request)
