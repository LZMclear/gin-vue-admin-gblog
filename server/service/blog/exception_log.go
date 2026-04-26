package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
)

type ExceptionLogService struct{}

func (s *ExceptionLogService) GetList(info blogReq.DateRangePageQuery) (list []blogModel.ExceptionLog, total int64, err error) {
	db := global.GVA_DB.Model(&blogModel.ExceptionLog{})
	if info.StartDate != "" && info.EndDate != "" {
		db = db.Where("create_time BETWEEN ? AND ?", info.StartDate, info.EndDate)
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

func (s *ExceptionLogService) Delete(id uint) error {
	return global.GVA_DB.Delete(&blogModel.ExceptionLog{}, id).Error
}
