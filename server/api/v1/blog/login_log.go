package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoginLogApi struct{}

func (a *LoginLogApi) GetLoginLogs(c *gin.Context) {
	var req blogReq.DateRangePageQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := loginLogService.GetList(req)
	if err != nil {
		global.GVA_LOG.Error("get login logs failed", zap.Error(err))
		response.FailWithMessage("获取登录日志失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: req.Page, PageSize: req.PageSize}, "获取登录日志成功", c)
}

func (a *LoginLogApi) DeleteLoginLog(c *gin.Context) {
	var req blogReq.IDReq
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := loginLogService.Delete(req.ID); err != nil {
		global.GVA_LOG.Error("delete login log failed", zap.Error(err))
		response.FailWithMessage("删除登录日志失败", c)
		return
	}
	response.OkWithMessage("删除登录日志成功", c)
}
