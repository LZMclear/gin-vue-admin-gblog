package request

import "time"

type MomentUpsert struct {
	ID          uint       `json:"id"`
	Content     string     `json:"content"`
	CreateTime  *time.Time `json:"createTime"`
	Likes       *int       `json:"likes"`
	IsPublished bool       `json:"isPublished"`
}
