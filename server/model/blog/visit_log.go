package blog

import "time"

type VisitLog struct {
	ID         uint      `json:"id" gorm:"column:id;primaryKey"`
	UUID       *string   `json:"uuid,omitempty" gorm:"column:uuid;size:36"`
	URI        string    `json:"uri" gorm:"column:uri;size:255;not null"`
	Method     string    `json:"method" gorm:"column:method;size:255;not null"`
	Param      string    `json:"param" gorm:"column:param;size:2000;not null"`
	Behavior   *string   `json:"behavior,omitempty" gorm:"column:behavior;size:255"`
	Content    *string   `json:"content,omitempty" gorm:"column:content;size:255"`
	Remark     *string   `json:"remark,omitempty" gorm:"column:remark;size:255"`
	IP         *string   `json:"ip,omitempty" gorm:"column:ip;size:255"`
	IPSource   *string   `json:"ipSource,omitempty" gorm:"column:ip_source;size:255"`
	OS         *string   `json:"os,omitempty" gorm:"column:os;size:255"`
	Browser    *string   `json:"browser,omitempty" gorm:"column:browser;size:255"`
	Times      int       `json:"times" gorm:"column:times;not null"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;not null"`
	UserAgent  *string   `json:"userAgent,omitempty" gorm:"column:user_agent;size:2000"`
}

func (VisitLog) TableName() string {
	return "gvto_visit_log"
}
