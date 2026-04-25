package blog

type Category struct {
	ID           uint   `json:"id" gorm:"column:id;primaryKey"`
	CategoryName string `json:"categoryName" gorm:"column:category_name;size:255;not null"`
}

func (Category) TableName() string {
	return "gvto_category"
}
