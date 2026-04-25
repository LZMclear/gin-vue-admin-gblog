package blog

type BlogTag struct {
	BlogID uint `json:"blogId" gorm:"column:blog_id;primaryKey"`
	TagID  uint `json:"tagId" gorm:"column:tag_id;primaryKey"`
}

func (BlogTag) TableName() string {
	return "gvto_blog_tag"
}
