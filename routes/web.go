package routes

import (
	"goblog/app/http/controllers"
	"goblog/app/http/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterWebRoutes 路由注册方法
func RegisterWebRoutes(r *mux.Router) {

	/**静态页面**/
	pc := new(controllers.PagesController)

	// tips gorilla/mux 使用的是精准匹配 如果路由不是 /about 而是 /about/ 多了一个斜杠的话也是找不处理方法的
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	// 自定义404页面
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
	/**文章相关页面**/
	ac := new(controllers.ArticlesController)
	r.HandleFunc("/", ac.Index).Methods("GET").Name("home")
	r.HandleFunc("/articles/{id:[0-9]+}", middlewares.Auth(ac.Show)).Methods("GET").Name("articles.show")
	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
	r.HandleFunc("/articles/create", ac.Create).Methods("GET").Name("articles.create")
	r.HandleFunc("/articles", ac.Store).Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/{id:[0-9]+}/edit", ac.Edit).Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Update).Methods("POST").Name("articles.update")
	r.HandleFunc("/articles/{id:[0-9]+}/delete", ac.Delete).Methods("POST").Name("articles.delete")

	/**静态资源路由**/
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	/**用户认证**/
	uac := new(controllers.AuthController)
	r.HandleFunc("/auth/register", uac.Register).Methods("GET").Name("auth.register")
	r.HandleFunc("/auth/do-register", uac.DoRegister).Methods("POST").Name("auth.doregister")
	r.HandleFunc("/auth/login", uac.Login).Methods("GET").Name("auth.login")
	r.HandleFunc("/auth/do-login", uac.DoLogin).Methods("POST").Name("auth.dologin")
	r.HandleFunc("/auth/logout", uac.Logout).Methods("POST").Name("auth.logout")

	// 使用中间件
	r.Use(middlewares.StartSession)
	//r.Use(middlewares.ForceHTML)
}
