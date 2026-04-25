package blog

import "time"

type Moment struct {
	ID          uint      `json:"id" gorm:"column:id;primaryKey"`
	Content     string    `json:"content" gorm:"column:content;type:longtext;not null"`
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time;not null"`
	Likes       *int      `json:"likes,omitempty" gorm:"column:likes"`
	IsPublished bool      `json:"isPublished" gorm:"column:is_published;not null"`
}

func (Moment) TableName() string {
	return "gvto_moment"
}
