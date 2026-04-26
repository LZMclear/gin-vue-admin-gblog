package blog

import "github.com/gin-gonic/gin"

type DashboardRouter struct{}

func (r *DashboardRouter) InitDashboardRouter(Router *gin.RouterGroup) {
	Router.Group("admin").GET("dashboard", dashboardApi.GetDashboard)
}
