package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TagRouter struct{}

func (r *TagRouter) InitTagRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	adminRecordRouter := Router.Group("admin").Use(middleware.OperationRecord())
	adminRouter := Router.Group("admin")
	{
		PublicRouter.GET("tag", tagApi.GetTagList)
	}
	{
		adminRouter.GET("tags", tagApi.GetTagList)
	}
	{
		adminRecordRouter.POST("tag", tagApi.CreateTag)
		adminRecordRouter.PUT("tag", tagApi.UpdateTag)
		adminRecordRouter.DELETE("tag", tagApi.DeleteTag)
	}
}
