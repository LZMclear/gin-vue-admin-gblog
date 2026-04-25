package blog

import "github.com/gin-gonic/gin"

type TelegramRouter struct{}

func (r *TelegramRouter) InitTelegramRouter(PublicRouter *gin.RouterGroup) {
	PublicRouter.POST("tg/:token", telegramApi.Webhook)
}
