package blog

import "time"

type ScheduleJobLog struct {
	LogID      uint       `json:"logId" gorm:"column:log_id;primaryKey;autoIncrement"`
	JobID      uint       `json:"jobId" gorm:"column:job_id;not null"`
	BeanName   *string    `json:"beanName,omitempty" gorm:"column:bean_name;size:255"`
	MethodName *string    `json:"methodName,omitempty" gorm:"column:method_name;size:255"`
	Params     *string    `json:"params,omitempty" gorm:"column:params;size:255"`
	Status     int8       `json:"status" gorm:"column:status;not null"`
	Error      *string    `json:"error,omitempty" gorm:"column:error;type:text"`
	Times      int        `json:"times" gorm:"column:times;not null"`
	CreateTime *time.Time `json:"createTime,omitempty" gorm:"column:create_time"`
}

func (ScheduleJobLog) TableName() string {
	return "gvto_schedule_job_log"
}
