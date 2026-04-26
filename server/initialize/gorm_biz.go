package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(
		blogModel.About{},
		blogModel.Blog{},
		blogModel.BlogTag{},
		blogModel.Category{},
		blogModel.CityVisitor{},
		blogModel.Comment{},
		blogModel.ExceptionLog{},
		blogModel.Friend{},
		blogModel.Moment{},
		blogModel.OperationLog{},
		blogModel.ScheduleJob{},
		blogModel.ScheduleJobLog{},
		blogModel.SiteSetting{},
		blogModel.Tag{},
		blogModel.User{},
		blogModel.VisitLog{},
		blogModel.VisitRecord{},
		blogModel.Visitor{},
	)
	if err != nil {
		return err
	}
	return nil
}
