package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
)

type AboutService struct{}

func (s *AboutService) GetList() ([]blogModel.About, error) {
	var list []blogModel.About
	err := global.GVA_DB.Order("id asc").Find(&list).Error
	return list, err
}

func (s *AboutService) UpdateValues(values map[string]string) error {
	tx := global.GVA_DB.Begin()
	for key, value := range values {
		if err := tx.Model(&blogModel.About{}).Where("name_en = ?", key).Update("value", value).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
