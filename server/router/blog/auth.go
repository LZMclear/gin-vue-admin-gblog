package blog

import (
	blogmw "github.com/flipped-aurora/gin-vue-admin/server/middleware/blog"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct{}

func (r *AuthRouter) InitAuthRouter(PublicRouter *gin.RouterGroup) {
	PublicRouter.Group("").Use(blogmw.VisitRecord(blogmw.VisitBehaviorCheckPassword)).POST("checkBlogPassword", authApi.CheckBlogPassword)
}
