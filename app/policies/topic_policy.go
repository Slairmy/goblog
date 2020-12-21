package policies

import (
	"goblog/app/models/article"
	"goblog/pkg/auth"
)

// CanModifyArticle 修改权限
func CanModifyArticle(_article article.Article) bool {
	return auth.User().ID == _article.UserID
}
