package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"html/template"
	"net/http"
	"strconv"
	"unicode/utf8"

	"gorm.io/gorm"
)

// ArticlesController 相当于PHP中定义Class
type ArticlesController struct {
}

// ArticlesFormData 给模版文件解析的变量结构体
type ArticlesFormData struct {
	Title, Content string
	URL            string
	Errors         map[string]string
}

// Index blog列表
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article.GetAll()

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "服务器错误")
	} else {
		// 加载模版
		tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		logger.LogError(err)
		tmpl.Execute(w, articles)
	}
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
	} else {
		fmt.Println(types.Int64ToString(article.ID))
		fmt.Println(route.Name2URL("articles.delete", "id", "2"))
		//tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
		// 高级用法 使用New先初始化,再注册函数,最后解析文件
		tmpl, err := template.New("show.gohtml").Funcs(
			template.FuncMap{
				"RouteName2URL": route.Name2URL,
				"Int64ToString": types.Int64ToString,
			}).ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)
		tmpl.Execute(w, article)
	}
}

// Create 创建blog
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	storeURL := route.Name2URL("articles.store")
	data := ArticlesFormData{
		Title:   "",
		Content: "",
		URL:     storeURL,
		Errors:  nil,
	}

	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, data)
}

// Store 创建blog
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	title := r.FormValue("title")
	content := r.FormValue("body")

	errors := validateArticleFormData(title, content)

	if len(errors) == 0 {
		// 验证成功保存
		_article := article.Article{
			Title:   title,
			Content: content,
		}
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功,ID为:"+strconv.FormatInt(_article.ID, 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}

	} else {

		storeURL := route.Name2URL("articles.store")

		data := ArticlesFormData{
			Title:   title,
			Content: content,
			URL:     storeURL,
			Errors:  errors,
		}

		tmpl, _ := template.ParseFiles("resources/views/articles/create.gohtml")
		tmpl.Execute(w, data)
	}
}

// Edit 创建blog
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章没找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器错误")
		}
	} else {

		// 拼接参数
		updateURL := route.Name2URL("articles.update", "id", id)

		logger.LogError(err)
		data := ArticlesFormData{
			Title:   article.Title,
			Content: article.Content,
			URL:     updateURL,
			Errors:  nil,
		}
		tmpl, _ := template.ParseFiles("resources/views/articles/edit.gohtml")
		tmpl.Execute(w, data)
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
		title := r.FormValue("title")
		content := r.FormValue("content")

		errors := validateArticleFormData(title, content)
		if len(errors) == 0 { // 无错误

			_article.Title = title
			_article.Content = content

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
			storeURL := route.Name2URL("articles.update", "id", id)
			data := ArticlesFormData{
				Title:   title,
				Content: content,
				URL:     storeURL,
				Errors:  errors,
			}

			tmpl, _ := template.ParseFiles("resources/views/articles/edit.gohtml")
			tmpl.Execute(w, data)
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
