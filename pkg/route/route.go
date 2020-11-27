package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router 全局路由变量
var Router *mux.Router

// Initialize 包初始化调用
func Initialize() {
	Router = mux.NewRouter()
}

// Name2URL 路由名称转换成URL路径
func Name2URL(routeName string, pairs ...string) string {
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		return ""
	}

	return url.String()
}

// GetRouteVariable 获取路由参数
func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
