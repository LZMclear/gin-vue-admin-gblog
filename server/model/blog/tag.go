package blog

type Tag struct {
	ID      uint    `json:"id" gorm:"column:id;primaryKey"`
	TagName string  `json:"tagName" gorm:"column:tag_name;size:255;not null"`
	Color   *string `json:"color,omitempty" gorm:"column:color;size:255"`
}

func (Tag) TableName() string {
	return "gvto_tag"
}
