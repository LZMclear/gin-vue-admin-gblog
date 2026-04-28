package blog

import "time"

type Comment struct {
	ID              uint       `json:"id" gorm:"column:id;primaryKey"`
	Nickname        string     `json:"nickname" gorm:"column:nickname;size:255;not null"`
	Email           string     `json:"email" gorm:"column:email;size:255;not null"`
	Content         string     `json:"content" gorm:"column:content;size:255;not null"`
	Avatar          string     `json:"avatar" gorm:"column:avatar;size:255;not null"`
	CreateTime      *time.Time `json:"createTime,omitempty" gorm:"column:create_time"`
	IP              *string    `json:"ip,omitempty" gorm:"column:ip;size:255"`
	IsPublished     bool       `json:"isPublished" gorm:"column:is_published;not null"`
	IsAdminComment  bool       `json:"isAdminComment" gorm:"column:is_admin_comment;not null"`
	Page            int        `json:"page" gorm:"column:page;not null"`
	IsNotice        bool       `json:"isNotice" gorm:"column:is_notice;not null"`
	BlogID          *uint      `json:"blogId,omitempty" gorm:"column:blog_id"`
	ParentCommentID int64      `json:"parentCommentId" gorm:"column:parent_comment_id;not null"`
	Website         *string    `json:"website,omitempty" gorm:"column:website;size:255"`
	QQ              *string    `json:"qq,omitempty" gorm:"column:qq;size:255"`
	Blog            *Blog      `json:"blog,omitempty" gorm:"-"`
}

func (Comment) TableName() string {
	return "gvto_comment"
}
