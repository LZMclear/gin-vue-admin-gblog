package blog

import "github.com/flipped-aurora/gin-vue-admin/server/model/common"

type TelegramService struct{}

func (s *TelegramService) HandleWebhook(token string, payload common.JSONMap) map[string]interface{} {
	return map[string]interface{}{
		"token":   token,
		"message": "telegram webhook received",
		"payload": payload,
	}
}
