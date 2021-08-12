package category

import "goblog/app/models"

type Category struct {
	models.BaseModel

	CategoryName string `gorm:"type:varchar(50);not null" valid:"categoryName"`
}
