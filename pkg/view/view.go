package view

import (
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// Render 模版渲染
func Render(w io.Writer, data interface{}, tplFiles ...string) {
	// 设置模版相对路径
	viewDir := "resources/views/"

	// 遍历所有的文件列表
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}

	files, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)

	allFiles := append(files, tplFiles...)

	tmpl, err := template.New("").Funcs(
		template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(allFiles...)
	logger.LogError(err)

	tmpl.ExecuteTemplate(w, "app", data)
}
