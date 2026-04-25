package blog

import "time"

type ExceptionLog struct {
	ID          uint      `json:"id" gorm:"column:id;primaryKey"`
	URI         string    `json:"uri" gorm:"column:uri;size:255;not null"`
	Method      string    `json:"method" gorm:"column:method;size:255;not null"`
	Param       *string   `json:"param,omitempty" gorm:"column:param;size:2000"`
	Description *string   `json:"description,omitempty" gorm:"column:description;size:255"`
	Error       *string   `json:"error,omitempty" gorm:"column:error;type:text"`
	IP          *string   `json:"ip,omitempty" gorm:"column:ip;size:255"`
	IPSource    *string   `json:"ipSource,omitempty" gorm:"column:ip_source;size:255"`
	OS          *string   `json:"os,omitempty" gorm:"column:os;size:255"`
	Browser     *string   `json:"browser,omitempty" gorm:"column:browser;size:255"`
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time;not null"`
	UserAgent   *string   `json:"userAgent,omitempty" gorm:"column:user_agent;size:2000"`
}

func (ExceptionLog) TableName() string {
	return "gvto_exception_log"
}
