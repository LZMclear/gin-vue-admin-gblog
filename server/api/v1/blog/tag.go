package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TagApi struct{}

func (a *TagApi) GetTagList(c *gin.Context) {
	list, err := tagService.GetList()
	if err != nil {
		global.GVA_LOG.Error("get tag list failed", zap.Error(err))
		response.FailWithMessage("获取标签列表失败", c)
		return
	}
	response.OkWithData(list, c)
}

func (a *TagApi) CreateTag(c *gin.Context) {
	var req blogReq.TagUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := tagService.Create(req); err != nil {
		global.GVA_LOG.Error("create tag failed", zap.Error(err))
		response.FailWithMessage("创建标签失败", c)
		return
	}
	response.OkWithMessage("创建标签成功", c)
}

func (a *TagApi) UpdateTag(c *gin.Context) {
	var req blogReq.TagUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := tagService.Update(req); err != nil {
		global.GVA_LOG.Error("update tag failed", zap.Error(err))
		response.FailWithMessage("更新标签失败", c)
		return
	}
	response.OkWithMessage("更新标签成功", c)
}

func (a *TagApi) DeleteTag(c *gin.Context) {
	var req blogReq.IDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := tagService.Delete(req.ID); err != nil {
		global.GVA_LOG.Error("delete tag failed", zap.Error(err))
		response.FailWithMessage("删除标签失败", c)
		return
	}
	response.OkWithMessage("删除标签成功", c)
}
