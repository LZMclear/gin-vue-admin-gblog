package blog

import (
	blogmw "github.com/flipped-aurora/gin-vue-admin/server/middleware/blog"
	"github.com/gin-gonic/gin"
)

type ArticleRouter struct{}

func (r *ArticleRouter) InitArticleRouter(PublicRouter *gin.RouterGroup) {
	PublicRouter.Group("").Use(blogmw.VisitRecord("index")).GET("blogs", articleApi.GetArticleList)
	PublicRouter.Group("").Use(blogmw.VisitRecord("blog")).GET("blog", articleApi.GetArticle)
}
