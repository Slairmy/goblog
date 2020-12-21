package controllers

import (
	"fmt"
	"goblog/pkg/flash"
	"goblog/pkg/logger"
	"net/http"

	"gorm.io/gorm"
)

type BaseController struct {
}

// ResponseForSQLError 统一响应
func (bc *BaseController) ResponseForSQLError(w http.ResponseWriter, err error) {
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章没找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器错误")
		}
	}
}

// ResponseForUnauthorized 统一响应未授权
func (bc *BaseController) ResponseForUnauthorized(w http.ResponseWriter, r *http.Request) {
	flash.Warning("没有权限操作")
	http.Redirect(w, r, "/", http.StatusForbidden)
}
