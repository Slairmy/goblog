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

	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

// GetAll 获取所有的blog
func GetAll() ([]Article, error) {
	var articles []Article
	if err := model.DB.Find(&articles).Error; err != nil {
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
