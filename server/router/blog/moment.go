package blog

import (
	blogmw "github.com/flipped-aurora/gin-vue-admin/server/middleware/blog"
	"github.com/gin-gonic/gin"
)

type MomentRouter struct{}

func (r *MomentRouter) InitMomentRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	adminRecordRouter := Router.Group("admin").Use(blogmw.OperationRecord())
	adminRouter := Router.Group("admin")
	PublicRouter.Group("").Use(blogmw.VisitRecord(blogmw.VisitBehaviorMoment)).GET("moments", momentApi.GetMoments)
	PublicRouter.Group("").Use(blogmw.VisitRecord(blogmw.VisitBehaviorLikeMoment)).POST("moment/like/:id", momentApi.LikeMoment)
	adminRouter.GET("moments", momentApi.GetMomentList)
	adminRouter.GET("moment", momentApi.GetMoment)
	adminRecordRouter.PUT("moment/published", momentApi.UpdateMomentPublished)
	adminRecordRouter.POST("moment", momentApi.CreateMoment)
	adminRecordRouter.PUT("moment", momentApi.UpdateMoment)
	adminRecordRouter.DELETE("moment", momentApi.DeleteMoment)
}
