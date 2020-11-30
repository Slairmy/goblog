package article

import (
	"goblog/pkg/route"
	"strconv"
)

// Article blog 类
type Article struct {
	ID      int64
	Title   string
	Content string
}

// Link 生成文章链接
func (a *Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatInt(a.ID, 10))
}
