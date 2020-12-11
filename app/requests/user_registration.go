package requests

import (
	"errors"
	"fmt"
	"goblog/app/models/user"
	"goblog/pkg/model"
	"strings"

	"github.com/thedevsaddam/govalidator"
)

// 此方法会在初始化时执行
func init() {
	govalidator.AddCustomRule("not_exists", func(field, rule, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")

		tableName := rng[0]
		dbField := rng[1]
		val := value.(string)

		var count int64
		model.DB.Table(tableName).Where(dbField+" = ?", val).Count(&count)
		if count != 0 {
			if message != "" {
				return errors.New(message)
			}

			return fmt.Errorf("%v 已被占用", val)
		}
		return nil
	})
}

// ValidateRegistrationForm 注册用户表单验证
func ValidateRegistrationForm(data user.User) map[string][]string {
	// 表单规则
	rules := govalidator.MapData{
		"name":             []string{"required", "alpha_num", "between:3,20", "not_exists:users,name"},
		"email":            []string{"required", "min:4", "max:30", "email", "not_exists:users,email"},
		"password":         []string{"required", "min:6"},
		"password_comfirm": []string{"required"},
	}

	// 定制错误信息
	message := govalidator.MapData{
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:格式错误,只允许英文和数字",
			"between:用户名长度需在3~20之间",
		},
		"email": []string{
			"required:Email 为必填项",
			"min:Email 长度需大于 4",
			"max:Email 长度需小于 30",
			"email:Email 格式不正确，请提供有效的邮箱地址",
		},
		"password": []string{
			"required:密码为必填项",
			"min:长度需大于 6",
		},
		"password_comfirm": []string{
			"required:确认密码框为必填项",
		},
	}

	// 表单验证
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid", // Struct 标签标识符
		Messages:      message,
	}
	errs := govalidator.New(opts).ValidateStruct()

	if data.Password != data.PasswordComfirm {
		errs["password_comfirm"] = append(errs["password_comfirm"], "两次输入密码不匹配")
	}

	return errs
}
