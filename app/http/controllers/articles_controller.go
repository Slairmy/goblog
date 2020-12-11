package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/requests"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
	"unicode/utf8"

	"gorm.io/gorm"
)

// ArticlesController 相当于PHP中定义Class
type ArticlesController struct {
}

// Index blog列表
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article.GetAll()

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "服务器错误")
	}

	view.Render(w, view.D{"Articles": articles}, "articles.index")
}

// Show 显示blog详情
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {

	id := route.GetRouteVariable("id", r)
	article, err := article.Get(id)

	if err != nil {
		// 返回空结果,但是没有报错
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章没找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器错误")
		}
	}
	// 渲染模版
	view.Render(w, view.D{"Article": article}, "articles.show")
}

// Create 创建blog
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "articles.create", "articles._form_field")
}

// Store 创建blog
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	title := r.FormValue("title")
	content := r.FormValue("content")

	_article := article.Article{
		Title:   title,
		Content: content,
	}
	errors := requests.ValidateArticleForm(_article)
	if len(errors) == 0 {
		// 验证成功保存
		_article.Create()
		if _article.ID > 0 {
			// 重定向一下
			http.Redirect(w, r, route.Name2URL("articles.index"), http.StatusFound)

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}

	} else {
		view.Render(w, view.D{
			"Article": _article,
			"Errors":  errors,
		}, "articles.create", "articles._form_field")
	}
}

// Edit 创建blog
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章没找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器错误")
		}
	} else {

		view.Render(w, view.D{
			"Article": _article,
		}, "articles.edit", "articles._form_field")

	}
}

// Update 创建blog
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章没找到")
		} else {
			// 打印日志
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器错误")
		}
	} else {
		// 无错误
		_article.Title = r.PostFormValue("titla")
		_article.Content = r.PostFormValue("content")

		errors := requests.ValidateArticleForm(_article)
		if len(errors) == 0 { // 无错误

			rowsAffected, err := _article.Update()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "服务器错误")
				return
			}

			if rowsAffected > 0 {
				showURL := route.Name2URL("articles.show", "id", id)
				http.Redirect(w, r, showURL, http.StatusFound)
			} else {
				fmt.Fprint(w, "无任何修改")
			}

		} else {
			view.Render(w, view.D{
				"Article": _article,
				"Errors":  errors,
			}, "articles.edit", "articles._form_field")
		}
	}
}

// Delete 删除blog
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章没找到")
		} else {
			// 打印日志
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器错误")
		}
	} else {
		rowAffected, err := _article.Delete()
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器错误")
		} else {
			if rowAffected > 0 {
				// 重定向
				indexURL := route.Name2URL("articles.index")
				http.Redirect(w, r, indexURL, http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章没找到")
			}
		}
	}
}

// 验证blog表单数据
func validateArticleFormData(title string, content string) map[string]string {
	errors := make(map[string]string)
	// tips: go 使用的utf-8编码格式,go 提供的len每个中文占3个字节,如果中文占一位的话,要使用utf8包
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 2 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度为2~40个字符"
	}

	if content == "" {
		errors["content"] = "内容不能为空"
	} else if utf8.RuneCountInString(content) < 10 {
		errors["content"] = "内容至少10个字符"
	}

	return errors
}
