package middlewares

import (
	"goblog/pkg/session"
	"net/http"
)

// StartSession 设置session中间件
func StartSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 开启session
		session.StartSession(w, r)

		next.ServeHTTP(w, r)
	})
}
