package blog

import "github.com/gin-gonic/gin"

type VisitorRouter struct{}

func (r *VisitorRouter) InitVisitorRouter(Router *gin.RouterGroup) {
	adminRouter := Router.Group("admin")
	adminRouter.GET("visitors", visitorApi.GetVisitors)
	adminRouter.DELETE("visitor", visitorApi.DeleteVisitor)
	adminRouter.DELETE("visitors", visitorApi.DeleteVisitors)
}
