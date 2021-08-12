package requests

import (
	"goblog/app/models/category"

	"github.com/thedevsaddam/govalidator"
)

func ValidateCategoryForm(data category.Category) map[string][]string {
	// 验证数据库 not_exists: table,column
	rules := govalidator.MapData{
		"categoryName": []string{"required", "min:2", "max:8", "not_exists:categories,category_name"},
	}

	messages := govalidator.MapData{
		"categoryName": []string{
			"required:请填写分类名称",
			"min:分类名称长度必须大于2",
			"max:分类名称长度必须小于8",
		},
	}

	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: "valid",
	}

	errs := govalidator.New(opts).ValidateStruct()

	return errs
}
