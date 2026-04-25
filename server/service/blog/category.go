package blog

import (
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
	entity := blogModel.Category{CategoryName: info.CategoryName}
	return global.GVA_DB.Create(&entity).Error
}

func (s *CategoryService) Update(info blogReq.CategoryUpsert) error {
	return global.GVA_DB.Model(&blogModel.Category{}).
		Where("id = ?", info.ID).
		Update("category_name", info.CategoryName).Error
}

func (s *CategoryService) Delete(id uint) error {
	return global.GVA_DB.Delete(&blogModel.Category{}, id).Error
}
