package blog

import (
	blogmw "github.com/flipped-aurora/gin-vue-admin/server/middleware/blog"
	"github.com/gin-gonic/gin"
)

type AboutRouter struct{}

func (r *AboutRouter) InitAboutRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	adminRecordRouter := Router.Group("admin").Use(blogmw.OperationRecord())
	adminRouter := Router.Group("admin")
	PublicRouter.Group("").Use(blogmw.VisitRecord("about")).GET("about", aboutApi.GetAbout)
	adminRouter.GET("about", aboutApi.GetAbout)
	adminRecordRouter.PUT("about", aboutApi.UpdateAbout)
}
