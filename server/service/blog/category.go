package blog

import (
	"errors"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
)

type CategoryService struct{}

func (s *CategoryService) GetList() ([]blogModel.Category, error) {
	var list []blogModel.Category
	err := global.GVA_DB.Order("id asc").Find(&list).Error
	return list, err
}

func (s *CategoryService) Create(info blogReq.CategoryUpsert) error {
	name := strings.TrimSpace(info.CategoryName)
	if name == "" {
		return errors.New("category name is required")
	}
	var existing blogModel.Category
	if err := global.GVA_DB.Where("category_name = ?", name).First(&existing).Error; err == nil {
		return errors.New("category already exists")
	}
	entity := blogModel.Category{CategoryName: name}
	return global.GVA_DB.Create(&entity).Error
}

func (s *CategoryService) Update(info blogReq.CategoryUpsert) error {
	name := strings.TrimSpace(info.CategoryName)
	if name == "" {
		return errors.New("category name is required")
	}
	var existing blogModel.Category
	if err := global.GVA_DB.Where("category_name = ?", name).First(&existing).Error; err == nil && existing.ID != info.ID {
		return errors.New("category already exists")
	}
	return global.GVA_DB.Model(&blogModel.Category{}).
		Where("id = ?", info.ID).
		Update("category_name", name).Error
}

func (s *CategoryService) Delete(id uint) error {
	var count int64
	if err := global.GVA_DB.Model(&blogModel.Blog{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("category is referenced by blogs")
	}
	return global.GVA_DB.Delete(&blogModel.Category{}, id).Error
}
