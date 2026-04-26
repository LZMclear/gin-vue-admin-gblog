package blog

import (
	blogmw "github.com/flipped-aurora/gin-vue-admin/server/middleware/blog"
	"github.com/gin-gonic/gin"
)

type ArchiveRouter struct{}

func (r *ArchiveRouter) InitArchiveRouter(PublicRouter *gin.RouterGroup) {
	PublicRouter.Group("").Use(blogmw.VisitRecord("archive")).GET("archives", archiveApi.GetArchives)
}
