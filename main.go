package main

import (
	"database/sql"
	"fmt"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

// ArticlesFormData 给模版文件解析的变量结构体
type ArticlesFormData struct {
	Title, Content string
	URL            *url.URL
	Errors         map[string]string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=mailto:summer@example.com>summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

// Article 动态
type Article struct {
	Title, Content string
	ID             int64
}

// Link 生成文章链接
func (a *Article) Link() string {
	showURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID, 10))
	if err != nil {
		logger.LogError(err)
		return ""
	}

	return showURL.String()
}

// Delete 删除文章
func (a *Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id = " + strconv.FormatInt(a.ID, 10))
	if err != nil {
		return 0, err
	}

	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}

	return 0, nil
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {

	id := route.GetRouteVariable("id", r)
	article, err := getArticleByID(id)

	if err != nil {
		// 返回空结果,但是没有报错
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章没找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器错误")
		}
	} else {
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

// RouteName2URL 路由名称转换成URL
func RouteName2URL(routeName string, pairs ...string) string {
	url, err := router.Get(routeName).URL(pairs...)
	if err != nil {
		logger.LogError(err)
		return ""
	}

	return url.String()
}

// Int64ToString 将 int64 转换为 string -- 分离types pkg
// func Int64ToString(num int64) string {
// 	return strconv.FormatInt(num, 10)
// }

// 文章列表
func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	// 返回结果集,包含数据库读取出来的数据和SQL连接
	rows, err := db.Query("SELECT * FROM articles")
	logger.LogError(err)
	// 注意这里是 row.Close()
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Content)
		logger.LogError(err)
		articles = append(articles, article)
	}

	// 加载模版
	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	logger.LogError(err)
	tmpl.Execute(w, articles)
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	// 解析表单数据
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "请提供正确的数据")
		return
	}

	// tips: go 使用的utf-8编码格式,go 提供的len每个中文占3个字节,如果中文占一位的话,要使用utf8包

	title := r.FormValue("title")
	content := r.FormValue("body")

	errors := validateArticleFormData(title, content)
	if len(errors) == 0 {
		// 验证成功保存
		lastInsertID, err := saveArticleToDB(title, content)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功,ID为:"+strconv.FormatInt(lastInsertID, 10))
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {

		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title:   title,
			Content: content,
			URL:     storeURL,
			Errors:  errors,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

		if err != nil {
			panic(err)
		}
		tmpl.Execute(w, data)
	}

}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
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

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {

	id := route.GetRouteVariable("id", r)
	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章没找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器错误")
		}
	} else {

		// 拼接参数
		updateURL, _ := router.Get("articles.update").URL("id", id)

		logger.LogError(err)
		data := ArticlesFormData{
			Title:   article.Title,
			Content: article.Content,
			URL:     updateURL,
			Errors:  nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")

		logger.LogError(err)
		tmpl.Execute(w, data)
	}
}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
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

			query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"

			// db.Exec是sql.DB方法 一般来处理 CREATE,UPDATE,DELETE方法 stmt.Exec 是sql.stmt 方法 使用Prepare 会发送两个sql请求
			rs, err := db.Exec(query, title, content, id)
			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "服务器错误")
			}

			// 更新成功
			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				// 重定向
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "没有任何修改")
			}

		} else {
			storeURL, _ := router.Get("articles.update").URL("id", id)

			data := ArticlesFormData{
				Title:   title,
				Content: content,
				URL:     storeURL,
				Errors:  errors,
			}

			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")

			if err != nil {
				panic(err)
			}
			tmpl.Execute(w, data)
		}
	}
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	article, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章没找到")
		} else {
			// 打印日志
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器错误")
		}
	} else {
		rowAffected, err := article.Delete()
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器错误")
		} else {
			if rowAffected > 0 {
				// 重定向
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章没找到")
			}
		}
	}
}

/**代码封装,目前获取根据文章id获取文章 和获取路由信息代码重复,可封装**/

// 获取路由参数 -- 移到route包
// func getRouteVariable(parameterName string, r *http.Request) string {
// 	vars := mux.Vars(r)
// 	return vars[parameterName]
// }

// 根据文章id获取文章信息 ? 为什么这里用string
func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"

	// 返回sql.Row 指针变量 保存有SQL连接,调用 Scan()会将连接释放
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Content)

	return article, err
}

// 验证blog表单数据
func validateArticleFormData(title string, content string) map[string]string {
	errors := make(map[string]string)

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

/**代码封装结束**/

// 中间件性质
func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// 中间件处理
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		next.ServeHTTP(w, r)
	})
}

// 分离 logger包记录错误日志
// func checkError(err error) {
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// 创建数据表
func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
		id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
		title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
		body longtext COLLATE utf8mb4_unicode_ci
	);`

	_, err := db.Exec(createArticlesSQL)
	logger.LogError(err)
}

// 保存文章
func saveArticleToDB(title string, body string) (int64, error) {
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	rs, err = stmt.Exec(title, body)

	if err != nil {
		return 0, err
	}

	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, nil
}

//var router = mux.NewRouter()
var router *mux.Router

// 全局db变量
var db *sql.DB

func main() {

	// 初始化路由
	route.Initialize()
	router = route.Router

	// 初始化DB
	database.Initialize()
	db = database.DB

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

	router.Use(forceHTMLMiddleware)

	homeURL, _ := router.Get("home").URL()
	fmt.Println("HomeUrl: ", homeURL)

	articleURL, _ := router.Get("articles.show").URL("id", "23")
	fmt.Println("articleURL: ", articleURL)

	// 文章详情
	router.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.SplitN(r.URL.Path, "/", 3)[2]
		fmt.Fprint(w, "文章ID: "+id)
	})

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
