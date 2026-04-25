package blog

import "time"

type LoginLog struct {
	ID          uint      `json:"id" gorm:"column:id;primaryKey"`
	Username    string    `json:"username" gorm:"column:username;size:255;not null"`
	IP          *string   `json:"ip,omitempty" gorm:"column:ip;size:255"`
	IPSource    *string   `json:"ipSource,omitempty" gorm:"column:ip_source;size:255"`
	OS          *string   `json:"os,omitempty" gorm:"column:os;size:255"`
	Browser     *string   `json:"browser,omitempty" gorm:"column:browser;size:255"`
	Status      *bool     `json:"status,omitempty" gorm:"column:status"`
	Description *string   `json:"description,omitempty" gorm:"column:description;size:255"`
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time;not null"`
	UserAgent   *string   `json:"userAgent,omitempty" gorm:"column:user_agent;size:2000"`
}

func (LoginLog) TableName() string {
	return "gvto_login_log"
}
