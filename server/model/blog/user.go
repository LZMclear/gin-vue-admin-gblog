package blog

import "time"

type User struct {
	ID         uint      `json:"id" gorm:"column:id;primaryKey"`
	Username   string    `json:"username" gorm:"column:username;size:255;not null"`
	Password   string    `json:"password" gorm:"column:password;size:255;not null"`
	Nickname   string    `json:"nickname" gorm:"column:nickname;size:255;not null"`
	Avatar     string    `json:"avatar" gorm:"column:avatar;size:255;not null"`
	Email      string    `json:"email" gorm:"column:email;size:255;not null"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;not null"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;not null"`
	Role       string    `json:"role" gorm:"column:role;size:255;not null"`
}

func (User) TableName() string {
	return "gvto_user"
}
