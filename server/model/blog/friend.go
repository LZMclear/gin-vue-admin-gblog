package blog

import "time"

type Friend struct {
	ID          uint      `json:"id" gorm:"column:id;primaryKey"`
	Nickname    string    `json:"nickname" gorm:"column:nickname;size:255;not null"`
	Description string    `json:"description" gorm:"column:description;size:255;not null"`
	Website     string    `json:"website" gorm:"column:website;size:255;not null"`
	Avatar      string    `json:"avatar" gorm:"column:avatar;size:255;not null"`
	IsPublished bool      `json:"isPublished" gorm:"column:is_published;not null"`
	Views       int       `json:"views" gorm:"column:views;not null"`
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time;not null"`
}

func (Friend) TableName() string {
	return "gvto_friend"
}
