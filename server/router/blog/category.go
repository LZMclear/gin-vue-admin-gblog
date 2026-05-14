package blog

import (
	blogmw "github.com/flipped-aurora/gin-vue-admin/server/middleware/blog"
	"github.com/gin-gonic/gin"
)

type CategoryRouter struct{}

func (r *CategoryRouter) InitCategoryRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	adminRecordRouter := Router.Group("admin").Use(blogmw.OperationRecord())
	adminRouter := Router.Group("admin")
	{
		PublicRouter.GET("category", categoryApi.GetCategoryList)
		PublicRouter.Group("").Use(blogmw.VisitRecord(blogmw.VisitBehaviorCategory)).GET("category/blogs", categoryApi.GetCategoryBlogList)
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
