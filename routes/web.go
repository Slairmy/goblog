package routes

import (
	"goblog/app/http/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterWebRoutes 路由注册方法
func RegisterWebRoutes(r *mux.Router) {

	/**静态页面**/
	pc := new(controllers.PagesController)

	// 路由文件
	r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	// tips gorilla/mux 使用的是精准匹配 如果路由不是 /about 而是 /about/ 多了一个斜杠的话也是找不处理方法的
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	// 自定义404页面
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
}
