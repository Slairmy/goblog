package article

import (
	"goblog/app/models"
	"goblog/pkg/route"
	"strconv"
)

// Article blog 类
type Article struct {
	models.BaseModel // 组合关系

	Title   string
	Content string
}

// Link 生成文章链接
func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatUint(a.ID, 10))
}
