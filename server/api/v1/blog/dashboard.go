package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DashboardApi struct{}

func (a *DashboardApi) GetDashboard(c *gin.Context) {
	data, err := dashboardService.GetSummary()
	if err != nil {
		global.GVA_LOG.Error("get dashboard failed", zap.Error(err))
		response.FailWithMessage("获取仪表盘数据失败", c)
		return
	}
	response.OkWithData(data, c)
}
