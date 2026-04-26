package blog

import "github.com/gin-gonic/gin"

type OperationLogRouter struct{}

func (r *OperationLogRouter) InitOperationLogRouter(Router *gin.RouterGroup) {
	adminRouter := Router.Group("admin")
	adminRouter.GET("operationLogs", operationLogApi.GetOperationLogs)
	adminRouter.DELETE("operationLog", operationLogApi.DeleteOperationLog)
}
