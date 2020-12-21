package article

import (
	"goblog/app/models"
	"goblog/app/models/user"
	"goblog/pkg/route"
	"strconv"
)

// Article blog 类
type Article struct {
	models.BaseModel // 组合关系

	Title   string `gorm:"type:varchar(255);not null" valid:"title"`
	Content string `gorm:"type:longtext;not null" valid:"content"`

	// 关联数据
	UserID uint64 `gorm:"not null;index"`
	User   user.User
}

// Link 生成文章链接
func (article Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatUint(article.ID, 10))
}

// CreatedAtDate 格式化时间
func (article Article) CreatedAtDate() string {
	return article.CreatedAt.Format("2006-01-02")
}
