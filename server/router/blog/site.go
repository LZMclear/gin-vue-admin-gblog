package blog

import "github.com/gin-gonic/gin"

type SiteRouter struct{}

func (r *SiteRouter) InitSiteRouter(Router *gin.RouterGroup) {
	Router.GET("site", siteApi.GetSite)
}
