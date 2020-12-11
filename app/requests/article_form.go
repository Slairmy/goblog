package requests

import (
	"goblog/app/models/article"

	"github.com/thedevsaddam/govalidator"
)

// ValidateArticleForm 验证表单
func ValidateArticleForm(data article.Article) map[string][]string {
	rules := govalidator.MapData{
		"title":   []string{"required", "min:3", "max:40"},
		"content": []string{"required", "min:10"},
	}

	messages := govalidator.MapData{
		"title": []string{
			"required:请填写标题",
			"min:标题长度必须大于3",
			"max:标题长度必须小于40",
		},
		"content": []string{
			"required:文章内容为必填项",
			"min:长度需大于 10",
		},
	}

	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      messages,
	}

	errs := govalidator.New(opts).ValidateStruct()

	return errs
}
