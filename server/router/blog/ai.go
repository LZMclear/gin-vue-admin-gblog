package blog

import "github.com/gin-gonic/gin"

type AiRouter struct{}

func (r *AiRouter) InitAiRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	aiRouter := Router.Group("blog/ai")
	{
		aiRouter.GET("status", aiApi.Status)
		aiRouter.POST("chat", aiApi.Chat)
		aiRouter.POST("summary", aiApi.Summary)
		aiRouter.POST("suggest-tags", aiApi.SuggestTags)
	}
}
