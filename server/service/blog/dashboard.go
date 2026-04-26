package blog

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
)

type DashboardService struct{}

func (s *DashboardService) GetSummary() (map[string]interface{}, error) {
	var blogCount int64
	if err := global.GVA_DB.Model(&blogModel.Blog{}).Count(&blogCount).Error; err != nil {
		return nil, err
	}

	var commentCount int64
	if err := global.GVA_DB.Model(&blogModel.Comment{}).Count(&commentCount).Error; err != nil {
		return nil, err
	}

	var todayPV int64
	if err := global.GVA_DB.Model(&blogModel.VisitLog{}).Where("date(create_time) = curdate()").Count(&todayPV).Error; err != nil {
		return nil, err
	}

	var todayUV int64
	if err := global.GVA_DB.Model(&blogModel.Visitor{}).Where("date(last_time) = curdate()").Count(&todayUV).Error; err != nil {
		return nil, err
	}

	var cityVisitors []blogModel.CityVisitor
	if err := global.GVA_DB.Order("uv desc").Limit(10).Find(&cityVisitors).Error; err != nil {
		return nil, err
	}

	var visitRecords []blogModel.VisitRecord
	if err := global.GVA_DB.Order("id desc").Limit(7).Find(&visitRecords).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"pv":           todayPV,
		"uv":           todayUV,
		"blogCount":    blogCount,
		"commentCount": commentCount,
		"category":     map[string]interface{}{"date": time.Now().Format("2006-01-02"), "list": []interface{}{}},
		"tag":          map[string]interface{}{"date": time.Now().Format("2006-01-02"), "list": []interface{}{}},
		"visitRecord":  visitRecords,
		"cityVisitor":  cityVisitors,
	}, nil
}
