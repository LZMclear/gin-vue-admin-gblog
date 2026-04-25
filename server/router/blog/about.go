package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AboutRouter struct{}

func (r *AboutRouter) InitAboutRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	adminRecordRouter := Router.Group("admin").Use(middleware.OperationRecord())
	adminRouter := Router.Group("admin")
	PublicRouter.GET("about", aboutApi.GetAbout)
	adminRouter.GET("about", aboutApi.GetAbout)
	adminRecordRouter.PUT("about", aboutApi.UpdateAbout)
}
