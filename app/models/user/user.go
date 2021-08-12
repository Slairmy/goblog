package user

import (
	"goblog/app/models"
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/password"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"time"
)

/** tips: GORM默认会将键小些转化为字段名称 **/

// User 用户信息
type User struct {
	models.BaseModel
	Name          string    `gorm:"type:varchar(255);not null;unique" valid:"name"`
	Avatar        string    `gorm:"type:varchar(255);default:NULL" valid:"avatar"`
	Introduction  string    `gorm:"type:varchar(255);default:NULL" valid:"introduction"`
	Email         string    `gorm:"type:varchar(255);default:NULL;unique" valid:"email"`
	Password      string    `gorm:"type:varchar(255)" valid:"password"`
	EmailVerifyAt time.Time `gorm:"column:email_verify_at"`

	// gorm: "-" 设置 GORM 在读写时忽略此字段
	PasswordComfirm string `gorm:"-" valid:"password_comfirm"`
}

// Link 生成用户链接
func (user User) Link() string {
	return route.Name2URL("users.show", "id", user.GetStringID())
}

// Create 创建用户
func (user *User) Create() (err error) {
	if err = model.DB.Create(&user).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

// GetByEmail 通过邮箱获取用户
func GetByEmail(email string) (User, error) {
	var user User
	if err := model.DB.Where("email = ?", email).Find(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// Get 通过session获取用户
func Get(id string) (User, error) {
	var user User
	uid := types.StringToInt(id)

	if err := model.DB.First(&user, uid).Error; err != nil {
		return user, err
	}

	return user, nil
}

// ComparePassword 比较密码
func (user *User) ComparePassword(_password string) bool {
	return password.CheckHash(_password, user.Password)
}
