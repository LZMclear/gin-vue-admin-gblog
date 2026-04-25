package blog

import "github.com/gin-gonic/gin"

type ExceptionLogRouter struct{}

func (r *ExceptionLogRouter) InitExceptionLogRouter(Router *gin.RouterGroup) {
	adminRouter := Router.Group("admin")
	adminRouter.GET("exceptionLogs", exceptionLogApi.GetExceptionLogs)
	adminRouter.DELETE("exceptionLog", exceptionLogApi.DeleteExceptionLog)
}
