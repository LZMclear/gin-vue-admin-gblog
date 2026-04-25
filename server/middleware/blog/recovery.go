package blog

import (
	"fmt"
	"net/http/httputil"
	"runtime/debug"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				req, _ := httputil.DumpRequest(c.Request, false)
				ip := c.ClientIP()
				ua := c.Request.UserAgent()
				desc := "blog panic"
				errMsg := fmt.Sprintf("panic: %v\nrequest: %s\nstack: %s", err, string(req), string(debug.Stack()))
				_ = global.GVA_DB.Create(&blogModel.ExceptionLog{
					URI:         c.Request.URL.Path,
					Method:      c.Request.Method,
					Param:       strPtr(limitString(c.Request.URL.RawQuery, 2000)),
					Description: &desc,
					Error:       &errMsg,
					IP:          &ip,
					CreateTime:  time.Now(),
					UserAgent:   &ua,
				}).Error
				panic(err)
			}
		}()
		c.Next()
	}
}
