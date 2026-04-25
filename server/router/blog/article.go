package blog

import "github.com/gin-gonic/gin"

type ArticleRouter struct{}

func (r *ArticleRouter) InitArticleRouter(PublicRouter *gin.RouterGroup) {
	PublicRouter.GET("blogs", articleApi.GetArticleList)
	PublicRouter.GET("blog", articleApi.GetArticle)
}
