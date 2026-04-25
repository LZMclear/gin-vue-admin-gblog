package blog

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	SiteApi         SiteApi
	ArchiveApi      ArchiveApi
	CategoryApi     CategoryApi
	TagApi          TagApi
	ArticleApi      ArticleApi
	AdminArticleApi AdminArticleApi
	AboutApi        AboutApi
	FriendApi       FriendApi
	MomentApi       MomentApi
	CommentApi      CommentApi
	SiteSettingApi  SiteSettingApi
	DashboardApi    DashboardApi
	VisitLogApi     VisitLogApi
	VisitorApi      VisitorApi
	ExceptionLogApi ExceptionLogApi
	LoginLogApi     LoginLogApi
	OperationLogApi OperationLogApi
	TelegramApi     TelegramApi
	AuthApi         AuthApi
}

var (
	siteService         = service.ServiceGroupApp.BlogServiceGroup.SiteService
	archiveService      = service.ServiceGroupApp.BlogServiceGroup.ArchiveService
	categoryService     = service.ServiceGroupApp.BlogServiceGroup.CategoryService
	tagService          = service.ServiceGroupApp.BlogServiceGroup.TagService
	articleService      = service.ServiceGroupApp.BlogServiceGroup.ArticleService
	adminArticleService = service.ServiceGroupApp.BlogServiceGroup.AdminArticleService
	aboutService        = service.ServiceGroupApp.BlogServiceGroup.AboutService
	friendService       = service.ServiceGroupApp.BlogServiceGroup.FriendService
	momentService       = service.ServiceGroupApp.BlogServiceGroup.MomentService
	commentService      = service.ServiceGroupApp.BlogServiceGroup.CommentService
	siteSettingService  = service.ServiceGroupApp.BlogServiceGroup.SiteSettingService
	dashboardService    = service.ServiceGroupApp.BlogServiceGroup.DashboardService
	visitLogService     = service.ServiceGroupApp.BlogServiceGroup.VisitLogService
	visitorService      = service.ServiceGroupApp.BlogServiceGroup.VisitorService
	exceptionLogService = service.ServiceGroupApp.BlogServiceGroup.ExceptionLogService
	loginLogService     = service.ServiceGroupApp.BlogServiceGroup.LoginLogService
	operationLogService = service.ServiceGroupApp.BlogServiceGroup.OperationLogService
	telegramService     = service.ServiceGroupApp.BlogServiceGroup.TelegramService
	authService         = service.ServiceGroupApp.BlogServiceGroup.AuthService
)
