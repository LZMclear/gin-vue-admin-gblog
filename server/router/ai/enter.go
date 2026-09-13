package ai

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	ModelConfigRouter
}

var modelConfigApi = api.ApiGroupApp.AiApiGroup.ModelConfigApi
