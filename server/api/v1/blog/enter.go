package blog

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	SiteApi         SiteApi
	CategoryApi     CategoryApi
	TagApi          TagApi
	ArticleApi      ArticleApi
	AdminArticleApi AdminArticleApi
}

var (
	siteService         = service.ServiceGroupApp.BlogServiceGroup.SiteService
	categoryService     = service.ServiceGroupApp.BlogServiceGroup.CategoryService
	tagService          = service.ServiceGroupApp.BlogServiceGroup.TagService
	articleService      = service.ServiceGroupApp.BlogServiceGroup.ArticleService
	adminArticleService = service.ServiceGroupApp.BlogServiceGroup.AdminArticleService
)
