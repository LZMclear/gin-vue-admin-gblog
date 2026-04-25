package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SiteApi struct{}

func (a *SiteApi) GetSite(c *gin.Context) {
	data, err := siteService.GetSiteInfo()
	if err != nil {
		global.GVA_LOG.Error("get site info failed", zap.Error(err))
		response.FailWithMessage("获取站点信息失败", c)
		return
	}
	response.OkWithData(data, c)
}
