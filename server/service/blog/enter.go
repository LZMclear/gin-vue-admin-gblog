package blog

type ServiceGroup struct {
	SiteService         SiteService
	ArchiveService      ArchiveService
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
	OperationLogService OperationLogCompatService
	TelegramService     TelegramService
	AuthService         AuthService
}
