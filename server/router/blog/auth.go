package blog

import "github.com/gin-gonic/gin"

type AuthRouter struct{}

func (r *AuthRouter) InitAuthRouter(PublicRouter *gin.RouterGroup) {
	PublicRouter.POST("checkBlogPassword", authApi.CheckBlogPassword)
}
