package blog

type CityVisitor struct {
	City string `json:"city" gorm:"column:city;size:255;primaryKey"`
	UV   int    `json:"uv" gorm:"column:uv;not null"`
}

func (CityVisitor) TableName() string {
	return "gvto_city_visitor"
}
