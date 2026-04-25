package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CategoryRouter struct{}

func (r *CategoryRouter) InitCategoryRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	adminRecordRouter := Router.Group("admin").Use(middleware.OperationRecord())
	adminRouter := Router.Group("admin")
	{
		PublicRouter.GET("category", categoryApi.GetCategoryList)
	}
	{
		adminRouter.GET("categories", categoryApi.GetCategoryList)
	}
	{
		adminRecordRouter.POST("category", categoryApi.CreateCategory)
		adminRecordRouter.PUT("category", categoryApi.UpdateCategory)
		adminRecordRouter.DELETE("category", categoryApi.DeleteCategory)
	}
}
