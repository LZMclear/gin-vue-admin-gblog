package blog

import (
	"bytes"
	"io"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
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
		ip := utils.BlogClientIP(c.Request)
		ipSource := utils.BlogIPSource(ip)
		ua := c.Request.UserAgent()
		os, browser := utils.BlogParseUserAgent(ua)
		dateKey := now.Format("2006-01-02")
		isNewVisitor := false
		var visitor blogModel.Visitor
		err := global.GVA_DB.Where("uuid = ?", identification).First(&visitor).Error
		if err == nil {
			_ = global.GVA_DB.Model(&blogModel.Visitor{}).Where("id = ?", visitor.ID).Updates(map[string]interface{}{
				"last_time": now,
				"pv":        gorm.Expr("coalesce(pv,0) + 1"),
				"ip":        ip,
				"ip_source": ipSource,
				"os":        os,
				"browser":   browser,
			}).Error
		} else if err == gorm.ErrRecordNotFound {
			isNewVisitor = true
			_ = global.GVA_DB.Create(&blogModel.Visitor{
				UUID:       identification,
				IP:         &ip,
				IPSource:   &ipSource,
				OS:         &os,
				Browser:    &browser,
				CreateTime: now,
				LastTime:   now,
				PV:         intPtr(1),
				UserAgent:  &ua,
			}).Error
			var city blogModel.CityVisitor
			if e := global.GVA_DB.Where("city = ?", ipSource).First(&city).Error; e == nil {
				_ = global.GVA_DB.Model(&blogModel.CityVisitor{}).Where("city = ?", ipSource).UpdateColumn("uv", gorm.Expr("uv + 1")).Error
			} else if e == gorm.ErrRecordNotFound {
				_ = global.GVA_DB.Create(&blogModel.CityVisitor{City: ipSource, UV: 1}).Error
			}
		}
		increaseDailyVisit(dateKey, isNewVisitor)

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
		_ = global.GVA_DB.Create(&blogModel.VisitLog{
			UUID:       &identification,
			URI:        c.Request.URL.Path,
			Method:     c.Request.Method,
			Param:      param,
			Behavior:   &behavior,
			IP:         &ip,
			IPSource:   &ipSource,
			OS:         &os,
			Browser:    &browser,
			Times:      int(cost.Milliseconds()),
			CreateTime: time.Now(),
			UserAgent:  &ua,
		}).Error
	}
}

func intPtr(v int) *int {
	return &v
}

func increaseDailyVisit(dateKey string, increaseUV bool) {
	var visitRecord blogModel.VisitRecord
	err := global.GVA_DB.Where("date = ?", dateKey).First(&visitRecord).Error
	if err == nil {
		updates := map[string]interface{}{
			"pv": gorm.Expr("pv + 1"),
		}
		if increaseUV {
			updates["uv"] = gorm.Expr("uv + 1")
		}
		_ = global.GVA_DB.Model(&blogModel.VisitRecord{}).Where("id = ?", visitRecord.ID).Updates(updates).Error
		return
	}
	if err == gorm.ErrRecordNotFound {
		record := blogModel.VisitRecord{Date: dateKey, PV: 1}
		if increaseUV {
			record.UV = 1
		}
		_ = global.GVA_DB.Create(&record).Error
	}
}
