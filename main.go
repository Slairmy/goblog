package main

import (
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"net/http"
)

func main() {

	// 初始化路由
	bootstrap.SetupDB()
	router := bootstrap.SetupRoute()

	http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
}
