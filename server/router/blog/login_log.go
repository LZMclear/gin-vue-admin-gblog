package blog

import "github.com/gin-gonic/gin"

type LoginLogRouter struct{}

func (r *LoginLogRouter) InitLoginLogRouter(Router *gin.RouterGroup) {
	adminRouter := Router.Group("admin")
	adminRouter.GET("loginLogs", loginLogApi.GetLoginLogs)
	adminRouter.DELETE("loginLog", loginLogApi.DeleteLoginLog)
}
