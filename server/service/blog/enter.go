package blog

type ServiceGroup struct {
	SiteService         SiteService
	CategoryService     CategoryService
	TagService          TagService
	ArticleService      ArticleService
	AdminArticleService AdminArticleService
}
