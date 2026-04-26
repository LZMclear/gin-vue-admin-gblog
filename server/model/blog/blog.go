package blog

import "time"

type Blog struct {
	ID               uint      `json:"id" gorm:"column:id;primaryKey"`
	Title            string    `json:"title" gorm:"column:title;size:255;not null"`
	FirstPicture     string    `json:"firstPicture" gorm:"column:first_picture;size:255;not null"`
	Content          string    `json:"content" gorm:"column:content;type:longtext;not null"`
	Description      string    `json:"description" gorm:"column:description;type:longtext;not null"`
	IsPublished      bool      `json:"isPublished" gorm:"column:is_published;not null"`
	IsRecommend      bool      `json:"isRecommend" gorm:"column:is_recommend;not null"`
	IsAppreciation   bool      `json:"isAppreciation" gorm:"column:is_appreciation;not null"`
	IsCommentEnabled bool      `json:"isCommentEnabled" gorm:"column:is_comment_enabled;not null"`
	CreateTime       time.Time `json:"createTime" gorm:"column:create_time;not null"`
	UpdateTime       time.Time `json:"updateTime" gorm:"column:update_time;not null"`
	Views            int       `json:"views" gorm:"column:views;not null"`
	Words            int       `json:"words" gorm:"column:words;not null"`
	ReadTime         int       `json:"readTime" gorm:"column:read_time;not null"`
	CategoryID       uint      `json:"categoryId" gorm:"column:category_id;not null"`
	IsTop            bool      `json:"isTop" gorm:"column:is_top;not null"`
	Password         *string   `json:"password,omitempty" gorm:"column:password;size:255"`
	UserID           *uint     `json:"userId,omitempty" gorm:"column:user_id"`
}

func (Blog) TableName() string {
	return "gvto_blog"
}
