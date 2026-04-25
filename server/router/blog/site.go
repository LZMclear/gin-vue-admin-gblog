package blog

import (
	blogmw "github.com/flipped-aurora/gin-vue-admin/server/middleware/blog"
	"github.com/gin-gonic/gin"
)

type SiteRouter struct{}

func (r *SiteRouter) InitSiteRouter(Router *gin.RouterGroup) {
	Router.Group("").Use(blogmw.VisitRecord("site")).GET("site", siteApi.GetSite)
}
