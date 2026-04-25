package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type TelegramApi struct{}

func (a *TelegramApi) Webhook(c *gin.Context) {
	token := c.Param("token")
	payload := common.JSONMap{}
	_ = c.ShouldBindJSON(&payload)
	response.OkWithData(telegramService.HandleWebhook(token, payload), c)
}
