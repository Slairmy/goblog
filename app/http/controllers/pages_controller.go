package controllers

import (
	"fmt"
	"net/http"
)

// PagesController 静态页面
type PagesController struct {
}

// Home 静态页面home处理
func (*PagesController) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

// About 静态页面关于页处理
func (*PagesController) About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=mailto:slairmy@163.com>slairmy@163.com</a>")
}

// NotFound 静态页面NotFound处理
func (*PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}
