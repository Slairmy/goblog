package controllers

import (
	"fmt"
	"goblog/app/models/category"
	"goblog/app/requests"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

// CategoryController 用户认证控制器
type CategoryController struct {
}

func (*CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "categories.create")
}

func (*CategoryController) Store(w http.ResponseWriter, r *http.Request) {
	categoryName := r.FormValue("categoryName")

	_category := category.Category{
		CategoryName: categoryName,
	}

	errors := requests.ValidateCategoryForm(_category)
	if len(errors) == 0 {
		_category.Create()
		if _category.ID > 0 {
			// 重定向到首页
			pageURL := route.Name2URL("home")
			http.Redirect(w, r, pageURL, http.StatusFound)
		} else {
			// 服务器错误
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		// 显示错误信息
		view.Render(w, view.D{
			"Errors": errors,
		}, "categories.create")
	}
}
