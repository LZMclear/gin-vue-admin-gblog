package blog

type VisitRecord struct {
	ID   uint   `json:"id" gorm:"column:id;primaryKey"`
	PV   int    `json:"pv" gorm:"column:pv;not null"`
	UV   int    `json:"uv" gorm:"column:uv;not null"`
	Date string `json:"date" gorm:"column:date;size:255;not null"`
}

func (VisitRecord) TableName() string {
	return "gvto_visit_record"
}
