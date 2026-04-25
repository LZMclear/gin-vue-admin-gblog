package blog

import "time"

type ScheduleJob struct {
	JobID      uint       `json:"jobId" gorm:"column:job_id;primaryKey"`
	BeanName   *string    `json:"beanName,omitempty" gorm:"column:bean_name;size:255"`
	MethodName *string    `json:"methodName,omitempty" gorm:"column:method_name;size:255"`
	Params     *string    `json:"params,omitempty" gorm:"column:params;size:255"`
	Cron       *string    `json:"cron,omitempty" gorm:"column:cron;size:255"`
	Status     *int8      `json:"status,omitempty" gorm:"column:status"`
	Remark     *string    `json:"remark,omitempty" gorm:"column:remark;size:255"`
	CreateTime *time.Time `json:"createTime,omitempty" gorm:"column:create_time"`
}

func (ScheduleJob) TableName() string {
	return "gvto_schedule_job"
}
