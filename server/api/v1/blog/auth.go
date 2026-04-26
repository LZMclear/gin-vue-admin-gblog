package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthApi struct{}

func (a *AuthApi) CheckBlogPassword(c *gin.Context) {
	var req blogReq.BlogPasswordCheck
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	token, err := authService.CreateBlogAccessToken(req.BlogID, req.Password)
	if err != nil {
		global.GVA_LOG.Error("check blog password failed", zap.Error(err))
		response.FailWithMessage("密码错误", c)
		return
	}
	response.OkWithData(token, c)
}
