package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type FriendRouter struct{}

func (r *FriendRouter) InitFriendRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	adminRecordRouter := Router.Group("admin").Use(middleware.OperationRecord())
	adminRouter := Router.Group("admin")
	PublicRouter.GET("friends", friendApi.GetFriends)
	PublicRouter.POST("friend", friendApi.AddFriendViews)
	adminRouter.GET("friends", friendApi.GetFriendList)
	adminRouter.GET("friendInfo", friendApi.GetFriendInfo)
	adminRecordRouter.PUT("friend/published", friendApi.UpdateFriendPublished)
	adminRecordRouter.POST("friend", friendApi.CreateFriend)
	adminRecordRouter.PUT("friend", friendApi.UpdateFriend)
	adminRecordRouter.DELETE("friend", friendApi.DeleteFriend)
	adminRecordRouter.PUT("friendInfo/commentEnabled", friendApi.UpdateFriendInfoCommentEnabled)
	adminRecordRouter.PUT("friendInfo/content", friendApi.UpdateFriendInfoContent)
}
