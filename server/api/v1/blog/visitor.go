package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VisitorApi struct{}

func (a *VisitorApi) GetVisitors(c *gin.Context) {
	var req blogReq.DateRangePageQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := visitorService.GetList(req)
	if err != nil {
		global.GVA_LOG.Error("get visitors failed", zap.Error(err))
		response.FailWithMessage("获取访客列表失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: req.Page, PageSize: req.PageSize}, "获取访客列表成功", c)
}

func (a *VisitorApi) DeleteVisitor(c *gin.Context) {
	var req blogReq.IDReq
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := visitorService.Delete(req.ID); err != nil {
		global.GVA_LOG.Error("delete visitor failed", zap.Error(err))
		response.FailWithMessage("删除访客失败", c)
		return
	}
	response.OkWithMessage("删除访客成功", c)
}
