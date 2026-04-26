package blog

import (
	"errors"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
)

type TagService struct{}

func (s *TagService) GetList() ([]blogModel.Tag, error) {
	var list []blogModel.Tag
	err := global.GVA_DB.Order("id asc").Find(&list).Error
	return list, err
}

func (s *TagService) Create(info blogReq.TagUpsert) error {
	name := strings.TrimSpace(info.TagName)
	if name == "" {
		return errors.New("tag name is required")
	}
	var existing blogModel.Tag
	if err := global.GVA_DB.Where("tag_name = ?", name).First(&existing).Error; err == nil {
		return errors.New("tag already exists")
	}
	entity := blogModel.Tag{TagName: name, Color: info.Color}
	return global.GVA_DB.Create(&entity).Error
}

func (s *TagService) Update(info blogReq.TagUpsert) error {
	name := strings.TrimSpace(info.TagName)
	if name == "" {
		return errors.New("tag name is required")
	}
	var existing blogModel.Tag
	if err := global.GVA_DB.Where("tag_name = ?", name).First(&existing).Error; err == nil && existing.ID != info.ID {
		return errors.New("tag already exists")
	}
	updates := map[string]interface{}{
		"tag_name": name,
		"color":    info.Color,
	}
	return global.GVA_DB.Model(&blogModel.Tag{}).
		Where("id = ?", info.ID).
		Updates(updates).Error
}

func (s *TagService) Delete(id uint) error {
	var count int64
	if err := global.GVA_DB.Model(&blogModel.BlogTag{}).Where("tag_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("tag is referenced by blogs")
	}
	return global.GVA_DB.Delete(&blogModel.Tag{}, id).Error
}
