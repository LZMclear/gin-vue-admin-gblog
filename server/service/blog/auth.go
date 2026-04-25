package blog

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
	blogResp "github.com/flipped-aurora/gin-vue-admin/server/model/blog/response"
	systemModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type AuthService struct{}

func (s *AuthService) Login(info blogReq.BlogLogin, ip, ua string) (blogResp.BlogLoginResponse, error) {
	user, err := systemService.UserServiceApp.Login(&systemModel.SysUser{
		Username: info.Username,
		Password: info.Password,
	})
	if err != nil {
		_ = s.writeLoginLogs(info.Username, ip, ua, false, "login failed")
		return blogResp.BlogLoginResponse{}, err
	}
	if user.Enable != 1 {
		_ = s.writeLoginLogs(info.Username, ip, ua, false, "user disabled")
		return blogResp.BlogLoginResponse{}, errors.New("user disabled")
	}
	token, err := utils.CreateBlogToken(utils.BlogAdminPrefix+user.Username, 7*24*time.Hour)
	if err != nil {
		return blogResp.BlogLoginResponse{}, err
	}
	_ = s.writeLoginLogsWithUser(user.ID, user.Username, ip, ua, true, "login success")
	return blogResp.BlogLoginResponse{
		User: systemRes.LoginResponse{
			User:      *user,
			Token:     token,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).UnixMilli(),
		},
		Token: token,
	}, nil
}

func (s *AuthService) CreateBlogAccessToken(blogID uint, password string) (string, error) {
	var blog blogModel.Blog
	if err := global.GVA_DB.Select("id,password").First(&blog, blogID).Error; err != nil {
		return "", err
	}
	if blog.Password == nil || strings.TrimSpace(*blog.Password) == "" {
		return "", errors.New("blog is not password protected")
	}
	if *blog.Password != password {
		return "", errors.New("password mismatch")
	}
	return utils.CreateBlogToken(utils.BlogAccessPrefix+uintToString(blogID), 30*24*time.Hour)
}

func (s *AuthService) writeLoginLogs(username, ip, ua string, status bool, desc string) error {
	return s.writeLoginLogsWithUser(0, username, ip, ua, status, desc)
}

func (s *AuthService) writeLoginLogsWithUser(userID uint, username, ip, ua string, status bool, desc string) error {
	ipSource := utils.BlogIPSource(ip)
	os, browser := utils.BlogParseUserAgent(ua)
	_ = systemService.LoginLogServiceApp.CreateLoginLog(systemModel.SysLoginLog{
		Username:     username,
		Ip:           ip,
		Status:       status,
		ErrorMessage: desc,
		Agent:        ua,
		UserID:       userID,
	})
	now := time.Now()
	return global.GVA_DB.Create(&blogModel.LoginLog{
		Username:    username,
		IP:          &ip,
		IPSource:    &ipSource,
		OS:          &os,
		Browser:     &browser,
		Status:      &status,
		Description: &desc,
		CreateTime:  now,
		UserAgent:   &ua,
	}).Error
}

func uintToString(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}
