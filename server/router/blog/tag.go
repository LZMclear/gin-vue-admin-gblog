package blog

import (
	blogmw "github.com/flipped-aurora/gin-vue-admin/server/middleware/blog"
	"github.com/gin-gonic/gin"
)

type TagRouter struct{}

func (r *TagRouter) InitTagRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	adminRecordRouter := Router.Group("admin").Use(blogmw.OperationRecord())
	adminRouter := Router.Group("admin")
	{
		PublicRouter.Group("").Use(blogmw.VisitRecord("tag")).GET("tag", tagApi.GetTagList)
		PublicRouter.Group("").Use(blogmw.VisitRecord("tag")).GET("tag/blogs", tagApi.GetTagBlogList)
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
