package blog

import "github.com/gin-gonic/gin"

type VisitLogRouter struct{}

func (r *VisitLogRouter) InitVisitLogRouter(Router *gin.RouterGroup) {
	adminRouter := Router.Group("admin")
	adminRouter.GET("visitLogs", visitLogApi.GetVisitLogs)
	adminRouter.DELETE("visitLog", visitLogApi.DeleteVisitLog)
	adminRouter.DELETE("visitLogs", visitLogApi.DeleteVisitLogs)
}
