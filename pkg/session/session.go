package session

import (
	"goblog/pkg/config"
	"goblog/pkg/logger"
	"net/http"

	"github.com/gorilla/sessions"
)

// 会话管理

// Store session容器
var Store = sessions.NewCookieStore([]byte(config.GetString("app.key")))

// Session 会话
var Session *sessions.Session

// Request 请求
var Request *http.Request

// Response 响应
var Response http.ResponseWriter

// StartSession 初始化session
func StartSession(w http.ResponseWriter, r *http.Request) {
	var err error
	Session, err = Store.Get(r, config.GetString("session.session_name"))

	logger.LogError(err)

	Request = r
	Response = w
}

/** session操作 **/

// Put session赋值
func Put(key string, value interface{}) {
	Session.Values[key] = value
	Save()

}

// Get 获取session值
func Get(key string) interface{} {
	return Session.Values[key]
}

// Forget 删除Session值
func Forget(key string) {
	delete(Session.Values, key)
	Save()
}

// Flush 刷新Session
func Flush() {
	Session.Options.MaxAge = -1
	Save()
}

// Save session存储
func Save() {
	err := Session.Save(Request, Response)
	logger.LogError(err)
}
