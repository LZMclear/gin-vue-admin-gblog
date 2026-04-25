package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AdminArticleRouter struct{}

func (r *AdminArticleRouter) InitAdminArticleRouter(Router *gin.RouterGroup) {
	adminRecordRouter := Router.Group("admin").Use(middleware.OperationRecord())
	adminRouter := Router.Group("admin")
	{
		adminRouter.GET("blogs", adminArticleApi.GetArticleList)
		adminRouter.GET("blog", adminArticleApi.GetArticle)
		adminRouter.GET("categoryAndTag", adminArticleApi.GetCategoryAndTag)
	}
	{
		adminRecordRouter.POST("blog", adminArticleApi.CreateArticle)
		adminRecordRouter.PUT("blog", adminArticleApi.UpdateArticle)
		adminRecordRouter.DELETE("blog", adminArticleApi.DeleteArticle)
	}
}
