package blog

import (
	"bytes"
	"io"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func VisitRecord(behavior string) gin.HandlerFunc {
	return func(c *gin.Context) {
		identification := c.GetHeader("identification")
		if identification == "" {
			identification = uuid.NewString()
			c.Header("identification", identification)
			c.Header("Access-Control-Expose-Headers", "identification")
		}

		now := time.Now()
		var visitor blogModel.Visitor
		err := global.GVA_DB.Where("uuid = ?", identification).First(&visitor).Error
		if err == nil {
			_ = global.GVA_DB.Model(&blogModel.Visitor{}).Where("id = ?", visitor.ID).Updates(map[string]interface{}{
				"last_time": now,
				"pv":        gorm.Expr("coalesce(pv,0) + 1"),
			}).Error
		} else if err == gorm.ErrRecordNotFound {
			ua := c.Request.UserAgent()
			ip := c.ClientIP()
			_ = global.GVA_DB.Create(&blogModel.Visitor{
				UUID:       identification,
				IP:         &ip,
				CreateTime: now,
				LastTime:   now,
				PV:         intPtr(1),
				UserAgent:  &ua,
			}).Error
		}

		var body []byte
		if c.Request.Method != "GET" && c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		start := time.Now()
		c.Next()
		cost := time.Since(start)

		param := c.Request.URL.RawQuery
		if len(body) > 0 {
			param = string(body)
		}
		ip := c.ClientIP()
		ua := c.Request.UserAgent()
		_ = global.GVA_DB.Create(&blogModel.VisitLog{
			UUID:       &identification,
			URI:        c.Request.URL.Path,
			Method:     c.Request.Method,
			Param:      param,
			Behavior:   &behavior,
			IP:         &ip,
			Times:      int(cost.Milliseconds()),
			CreateTime: time.Now(),
			UserAgent:  &ua,
		}).Error
	}
}

func intPtr(v int) *int {
	return &v
}
