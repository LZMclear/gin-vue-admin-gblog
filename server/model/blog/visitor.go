package blog

import "time"

type Visitor struct {
	ID         uint      `json:"id" gorm:"column:id;primaryKey"`
	UUID       string    `json:"uuid" gorm:"column:uuid;size:36;uniqueIndex:idx_uuid;not null"`
	IP         *string   `json:"ip,omitempty" gorm:"column:ip;size:255"`
	IPSource   *string   `json:"ipSource,omitempty" gorm:"column:ip_source;size:255"`
	OS         *string   `json:"os,omitempty" gorm:"column:os;size:255"`
	Browser    *string   `json:"browser,omitempty" gorm:"column:browser;size:255"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;not null"`
	LastTime   time.Time `json:"lastTime" gorm:"column:last_time;not null"`
	PV         *int      `json:"pv,omitempty" gorm:"column:pv"`
	UserAgent  *string   `json:"userAgent,omitempty" gorm:"column:user_agent;size:2000"`
}

func (Visitor) TableName() string {
	return "gvto_visitor"
}
