package blog

import (
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
	entity := blogModel.Tag{TagName: info.TagName, Color: info.Color}
	return global.GVA_DB.Create(&entity).Error
}

func (s *TagService) Update(info blogReq.TagUpsert) error {
	updates := map[string]interface{}{
		"tag_name": info.TagName,
		"color":    info.Color,
	}
	return global.GVA_DB.Model(&blogModel.Tag{}).
		Where("id = ?", info.ID).
		Updates(updates).Error
}

func (s *TagService) Delete(id uint) error {
	return global.GVA_DB.Delete(&blogModel.Tag{}, id).Error
}
