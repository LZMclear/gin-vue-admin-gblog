package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SiteSettingRouter struct{}

func (r *SiteSettingRouter) InitSiteSettingRouter(Router *gin.RouterGroup) {
	adminRecordRouter := Router.Group("admin").Use(middleware.OperationRecord())
	adminRouter := Router.Group("admin")
	adminRouter.GET("siteSettings", siteSettingApi.GetSiteSettings)
	adminRouter.GET("webTitleSuffix", siteSettingApi.GetWebTitleSuffix)
	adminRecordRouter.POST("siteSettings", siteSettingApi.UpdateSiteSettings)
}
