package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminArticleApi struct{}

func (a *AdminArticleApi) GetArticleList(c *gin.Context) {
	var req blogReq.AdminArticleSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := adminArticleService.GetList(req)
	if err != nil {
		global.GVA_LOG.Error("get admin article list failed", zap.Error(err))
		response.FailWithMessage("获取后台文章列表失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取后台文章列表成功", c)
}

func (a *AdminArticleApi) GetArticle(c *gin.Context) {
	var req blogReq.IDQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := adminArticleService.GetByID(req.ID)
	if err != nil {
		global.GVA_LOG.Error("get admin article failed", zap.Error(err))
		response.FailWithMessage("获取后台文章失败", c)
		return
	}
	response.OkWithData(data, c)
}

func (a *AdminArticleApi) GetCategoryAndTag(c *gin.Context) {
	data, err := adminArticleService.GetCategoryAndTag()
	if err != nil {
		global.GVA_LOG.Error("get category and tag failed", zap.Error(err))
		response.FailWithMessage("获取分类和标签失败", c)
		return
	}
	response.OkWithData(data, c)
}

func (a *AdminArticleApi) CreateArticle(c *gin.Context) {
	var req blogReq.ArticleUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := adminArticleService.Create(req); err != nil {
		global.GVA_LOG.Error("create article failed", zap.Error(err))
		response.FailWithMessage("创建文章失败", c)
		return
	}
	response.OkWithMessage("创建文章成功", c)
}

func (a *AdminArticleApi) UpdateArticle(c *gin.Context) {
	var req blogReq.ArticleUpsert
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := adminArticleService.Update(req); err != nil {
		global.GVA_LOG.Error("update article failed", zap.Error(err))
		response.FailWithMessage("更新文章失败", c)
		return
	}
	response.OkWithMessage("更新文章成功", c)
}

func (a *AdminArticleApi) DeleteArticle(c *gin.Context) {
	var req blogReq.IDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := adminArticleService.Delete(req.ID); err != nil {
		global.GVA_LOG.Error("delete article failed", zap.Error(err))
		response.FailWithMessage("删除文章失败", c)
		return
	}
	response.OkWithMessage("删除文章成功", c)
}
