package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ArticleApi struct{}

func (a *ArticleApi) GetArticleList(c *gin.Context) {
	var req blogReq.ArticleSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := articleService.GetPublishedList(req)
	if err != nil {
		global.GVA_LOG.Error("get article list failed", zap.Error(err))
		response.FailWithMessage("获取文章列表失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取文章列表成功", c)
}

func (a *ArticleApi) GetArticle(c *gin.Context) {
	var req blogReq.IDQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := articleService.GetPublishedByID(req.ID)
	if err != nil {
		global.GVA_LOG.Error("get article failed", zap.Error(err))
		response.FailWithMessage("获取文章失败", c)
		return
	}
	response.OkWithData(data, c)
}
