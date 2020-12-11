package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"goblog/pkg/session"
	"goblog/pkg/view"
	"net/http"
)

// AuthController 用户认证控制器
type AuthController struct {
}

type userForm struct {
	Name            string `valid:"name"`
	Email           string `valid:"email"`
	Password        string `valid:"password"`
	PasswordComfirm string `valid:"password_comfirm"`
}

// Register 注册逻辑
func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

// DoRegister 处理注册逻辑
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {

	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordComfirm: r.PostFormValue("password_confirmation"),
	}

	// 表单验证
	errs := requests.ValidateRegistrationForm(_user)

	if len(errs) > 0 {
		// data, _ := json.MarshalIndent(errs, "", " ")
		// fmt.Fprint(w, string(data))
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")
	} else {
		_user.Create()

		if _user.ID > 0 {
			flash.Success("恭喜您注册成功!")
			// 注册成功之后登陆
			auth.Login(_user)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "注册失败,请联系管理员")
		}
	}

	// 表单不通过重新显示表单

}

// Login 登陆界面
func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(session.Get("uid"))
	view.RenderSimple(w, view.D{}, "auth.login")
}

// DoLogin 处理登陆逻辑
func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	if err := auth.Attempt(email, password); err == nil {
		flash.Success("欢迎回来！")
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		view.RenderSimple(w, view.D{
			"Error":    err.Error(),
			"Email":    email,
			"Password": password,
		}, "auth.login")
	}
}

// Logout 退出登陆
func (*AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// 清除session
	auth.Logout()
	flash.Success("您已退出登陆！")
	http.Redirect(w, r, "/", http.StatusFound)
}
