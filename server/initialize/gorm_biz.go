package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	aiModel "github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(
		aiModel.AiModelConfig{},
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
		blogModel.ScheduleJobLog{},
		blogModel.SiteSetting{},
		blogModel.Tag{},
		blogModel.VisitLog{},
		blogModel.VisitRecord{},
		blogModel.Visitor{},
	)
	if err != nil {
		return err
	}
	return nil
}
