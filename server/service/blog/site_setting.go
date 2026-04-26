package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
	"gorm.io/gorm"
)

type SiteSettingService struct{}

func (s *SiteSettingService) GetGrouped() (map[string][]blogModel.SiteSetting, error) {
	var list []blogModel.SiteSetting
	if err := global.GVA_DB.Order("id asc").Find(&list).Error; err != nil {
		return nil, err
	}
	result := map[string][]blogModel.SiteSetting{}
	for _, item := range list {
		key := "type0"
		if item.Type != nil {
			key = "type" + string(rune('0'+*item.Type))
		}
		result[key] = append(result[key], item)
	}
	return result, nil
}

func (s *SiteSettingService) UpdateAll(info blogReq.SiteSettingBatchUpdate) error {
	tx := global.GVA_DB.Begin()
	for _, id := range info.DeleteIDs {
		if err := tx.Delete(&blogModel.SiteSetting{}, id).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, item := range info.Settings {
		if item.ID > 0 {
			if err := tx.Model(&blogModel.SiteSetting{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
				"name_en": item.NameEn,
				"name_zh": item.NameZh,
				"value":   item.Value,
				"type":    item.Type,
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			entity := blogModel.SiteSetting{
				NameEn: item.NameEn,
				NameZh: item.NameZh,
				Value:  item.Value,
				Type:   item.Type,
			}
			if err := tx.Create(&entity).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

func (s *SiteSettingService) GetWebTitleSuffix() (string, error) {
	var row blogModel.SiteSetting
	err := global.GVA_DB.Where("name_en = ?", "webTitleSuffix").First(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	if row.Value == nil {
		return "", nil
	}
	return *row.Value, nil
}
