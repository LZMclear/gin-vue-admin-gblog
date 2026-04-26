package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ArchiveApi struct{}

func (a *ArchiveApi) GetArchives(c *gin.Context) {
	var req blogReq.ArchiveSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := archiveService.GetArchive(req)
	if err != nil {
		global.GVA_LOG.Error("get archives failed", zap.Error(err))
		response.FailWithMessage("获取归档失败", c)
		return
	}
	response.OkWithData(data, c)
}
