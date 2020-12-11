package auth

import (
	"errors"
	"goblog/app/models/user"
	"goblog/pkg/session"

	"gorm.io/gorm"
)

func _getUID() string {
	_uid := session.Get("uid")

	// 类型断言: 因为go所有的类型都实现了interface{}接口,将interface{}转换成对应的类型需要断言
	uid, ok := _uid.(string)

	if ok && len(uid) > 0 {
		return uid
	}
	return ""
}

// User 通过session获取用户
func User() user.User {
	uid := _getUID()
	if len(uid) > 0 {
		_user, err := user.Get(uid)
		if err == nil {
			return _user
		}
	}

	return user.User{}
}

// Attempt 尝试登陆
func Attempt(email string, password string) error {
	_user, err := user.GetByEmail(email)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("账号不存在或密码错误")
		}

		return errors.New("内部错误,请稍后再试")
	}

	// 匹配密码
	if !_user.ComparePassword(password) {
		return errors.New("密码错误")
	}

	// 缓存的是ID
	session.Put("uid", _user.GetStringID())

	return nil
}

// Login 登陆
func Login(_user user.User) {
	session.Put("uid", _user.GetStringID())
}

// Logout 登出
func Logout() {
	session.Forget("uid")
}

// Check 检测cookie
func Check() bool {
	return len(_getUID()) > 0
}
