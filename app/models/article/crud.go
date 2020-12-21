package article

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/types"
)

// Get 根据id获取blog
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToInt(idstr)

	// 预加载关联数据
	if err := model.DB.Preload("User").First(&article, id).Error; err != nil {
		//if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

// GetAll 获取所有的blog
func GetAll() ([]Article, error) {
	var articles []Article

	// 预加载关联数据
	// 调试单条模型可以使用 GORM提供的Debug方法 model.DB.Debug().Preload(model).Find().Error
	if err := model.DB.Preload("User").Find(&articles).Error; err != nil {
		//if err := model.DB.Find(&articles).Error; err != nil {

		return articles, err
	}

	return articles, nil
}

// Create 插入blog
func (article *Article) Create() (err error) {
	if err = model.DB.Create(&article).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

// Update 更新blog
func (article *Article) Update() (rowsAffected int64, err error) {
	result := model.DB.Save(&article)
	if err = result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}

// Delete 删除blog
func (article *Article) Delete() (rowsAffected int64, err error) {
	result := model.DB.Delete(&article)
	if err = result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}

// GetByUserID 通过userID获取文章
func GetByUserID(userID string) ([]Article, error) {
	var articles []Article

	if err := model.DB.Preload("User").Where("user_id = ?", userID).Find(&articles).Error; err != nil {
		return articles, err
	}

	return articles, nil
}
