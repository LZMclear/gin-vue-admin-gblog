package blog

type ServiceGroup struct {
	SiteService         SiteService
	CategoryService     CategoryService
	TagService          TagService
	ArticleService      ArticleService
	AdminArticleService AdminArticleService
	AboutService        AboutService
	FriendService       FriendService
	MomentService       MomentService
	CommentService      CommentService
	SiteSettingService  SiteSettingService
	DashboardService    DashboardService
	VisitLogService     VisitLogService
	VisitorService      VisitorService
	ExceptionLogService ExceptionLogService
	LoginLogService     LoginLogCompatService
	OperationLogService OperationLogCompatService
	TelegramService     TelegramService
}
