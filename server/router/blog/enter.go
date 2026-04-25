package blog

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	SiteRouter
	CategoryRouter
	TagRouter
	ArticleRouter
	AdminArticleRouter
	AboutRouter
	FriendRouter
	MomentRouter
	CommentRouter
	SiteSettingRouter
	DashboardRouter
}

var (
	siteApi         = api.ApiGroupApp.BlogApiGroup.SiteApi
	categoryApi     = api.ApiGroupApp.BlogApiGroup.CategoryApi
	tagApi          = api.ApiGroupApp.BlogApiGroup.TagApi
	articleApi      = api.ApiGroupApp.BlogApiGroup.ArticleApi
	adminArticleApi = api.ApiGroupApp.BlogApiGroup.AdminArticleApi
	aboutApi        = api.ApiGroupApp.BlogApiGroup.AboutApi
	friendApi       = api.ApiGroupApp.BlogApiGroup.FriendApi
	momentApi       = api.ApiGroupApp.BlogApiGroup.MomentApi
	commentApi      = api.ApiGroupApp.BlogApiGroup.CommentApi
	siteSettingApi  = api.ApiGroupApp.BlogApiGroup.SiteSettingApi
	dashboardApi    = api.ApiGroupApp.BlogApiGroup.DashboardApi
)
