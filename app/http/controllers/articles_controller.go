package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/policies"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
	"unicode/utf8"
)

// ArticlesController 相当于PHP中定义Class
type ArticlesController struct {
	BaseController
}

// Index 动态列表
func (ac *ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article.GetAll()
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		view.Render(w, view.D{
			"Articles": articles,
		}, "articles.index", "articles._article_meta")
	}
}

// Show 动态详情
func (ac *ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)
	if err != nil {
		// 响应SQL
		ac.ResponseForSQLError(w, err)
	} else {
		// 渲染模版
		view.Render(w, view.D{
			"Article":          _article,
			"CanModifyArticle": policies.CanModifyArticle(_article),
		}, "articles.show", "articles._article_meta")
	}
}

// Create 创建动态
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "articles.create", "articles._form_field")
}

// Store 保存动态
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	title := r.FormValue("title")
	content := r.FormValue("content")

	_article := article.Article{
		Title:   title,
		Content: content,
		UserID:  auth.User().ID,
	}
	errors := requests.ValidateArticleForm(_article)
	if len(errors) == 0 {
		// 验证成功保存
		_article.Create()
		if _article.ID > 0 {
			// 重定向到创建的动态
			showURL := route.Name2URL("articles.show", "id", _article.GetStringID())
			http.Redirect(w, r, showURL, http.StatusFound)

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

// Edit 编辑动态
func (ac *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 查看权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			view.Render(w, view.D{
				"Article": _article,
			}, "articles.edit", "articles._form_field")
		}
	}
}

// Update 创建动态
func (ac *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			// 无错误
			_article.Title = r.PostFormValue("title")
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
}

// Delete 删除动态
func (ac *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		}
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
