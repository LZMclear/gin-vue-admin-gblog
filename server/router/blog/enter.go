package blog

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	SiteRouter
	CategoryRouter
	TagRouter
	ArticleRouter
	AdminArticleRouter
}

var (
	siteApi         = api.ApiGroupApp.BlogApiGroup.SiteApi
	categoryApi     = api.ApiGroupApp.BlogApiGroup.CategoryApi
	tagApi          = api.ApiGroupApp.BlogApiGroup.TagApi
	articleApi      = api.ApiGroupApp.BlogApiGroup.ArticleApi
	adminArticleApi = api.ApiGroupApp.BlogApiGroup.AdminArticleApi
)
