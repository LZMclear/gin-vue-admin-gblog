package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
)

type VisitorService struct{}

func (s *VisitorService) GetList(info blogReq.DateRangePageQuery) (list []blogModel.Visitor, total int64, err error) {
	db := global.GVA_DB.Model(&blogModel.Visitor{})
	if info.StartDate != "" && info.EndDate != "" {
		db = db.Where("last_time BETWEEN ? AND ?", info.StartDate, info.EndDate)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	offset := info.PageSize * (info.Page - 1)
	err = db.Order("create_time desc").Limit(info.PageSize).Offset(offset).Find(&list).Error
	return
}

func (s *VisitorService) Delete(id uint) error {
	return global.GVA_DB.Delete(&blogModel.Visitor{}, id).Error
}
