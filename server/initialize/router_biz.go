package initialize

import (
	blogmw "github.com/flipped-aurora/gin-vue-admin/server/middleware/blog"
	"github.com/flipped-aurora/gin-vue-admin/server/router"
	"github.com/gin-gonic/gin"
)

// 占位方法，保证文件可以正确加载，避免go空变量检测报错，请勿删除。
func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroupApp
}

func initBizRouter(routers ...*gin.RouterGroup) {
	privateGroup := routers[0].Group("")
	privateGroup.Use(blogmw.Recovery())
	publicGroup := routers[1].Group("")
	publicGroup.Use(blogmw.Recovery())
	blogRouter := router.RouterGroupApp.Blog

	holder(publicGroup, privateGroup)

	blogRouter.InitSiteRouter(publicGroup)
	blogRouter.InitAuthRouter(publicGroup)
	blogRouter.InitCategoryRouter(privateGroup, publicGroup)
	blogRouter.InitTagRouter(privateGroup, publicGroup)
	blogRouter.InitArticleRouter(publicGroup)
	blogRouter.InitAdminArticleRouter(privateGroup)
	blogRouter.InitAboutRouter(privateGroup, publicGroup)
	blogRouter.InitFriendRouter(privateGroup, publicGroup)
	blogRouter.InitMomentRouter(privateGroup, publicGroup)
	blogRouter.InitCommentRouter(privateGroup, publicGroup)
	blogRouter.InitSiteSettingRouter(privateGroup)
	blogRouter.InitDashboardRouter(privateGroup)
	blogRouter.InitVisitLogRouter(privateGroup)
	blogRouter.InitVisitorRouter(privateGroup)
	blogRouter.InitExceptionLogRouter(privateGroup)
	blogRouter.InitLoginLogRouter(privateGroup)
	blogRouter.InitOperationLogRouter(privateGroup)
	blogRouter.InitTelegramRouter(publicGroup)

}
