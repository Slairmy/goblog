package bootstrap

import (
	"goblog/pkg/route"
	"goblog/routes"

	"github.com/gorilla/mux"
)

// SetupRoute 项目开始初始化路由
func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)

	route.SetRoute(router)
	return router
}
