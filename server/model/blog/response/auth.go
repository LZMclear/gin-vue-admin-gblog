package response

import systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"

type BlogLoginResponse struct {
	User  systemRes.LoginResponse `json:"user"`
	Token string                  `json:"token"`
}
