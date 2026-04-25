package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type MomentRouter struct{}

func (r *MomentRouter) InitMomentRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	adminRecordRouter := Router.Group("admin").Use(middleware.OperationRecord())
	adminRouter := Router.Group("admin")
	PublicRouter.GET("moments", momentApi.GetMoments)
	PublicRouter.POST("moment/like/:id", momentApi.LikeMoment)
	adminRouter.GET("moments", momentApi.GetMomentList)
	adminRouter.GET("moment", momentApi.GetMoment)
	adminRecordRouter.PUT("moment/published", momentApi.UpdateMomentPublished)
	adminRecordRouter.POST("moment", momentApi.CreateMoment)
	adminRecordRouter.PUT("moment", momentApi.UpdateMoment)
	adminRecordRouter.DELETE("moment", momentApi.DeleteMoment)
}
