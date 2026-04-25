package blog

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
)

type LoginLogCompatService struct{}

func (s *LoginLogCompatService) GetList(info blogReq.DateRangePageQuery) (list []blogModel.LoginLog, total int64, err error) {
	db := global.GVA_DB.Model(&blogModel.LoginLog{})
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

func (s *LoginLogCompatService) Delete(id uint) error {
	return global.GVA_DB.Delete(&blogModel.LoginLog{}, id).Error
}

func (s *LoginLogCompatService) Create(username, ip, userAgent string, status bool, description string) error {
	now := time.Now()
	return global.GVA_DB.Create(&blogModel.LoginLog{
		Username:    username,
		IP:          &ip,
		Status:      &status,
		Description: &description,
		CreateTime:  now,
		UserAgent:   &userAgent,
	}).Error
}
