package blog

type About struct {
	ID     uint   `json:"id" gorm:"column:id;primaryKey"`
	NameEn string `json:"nameEn" gorm:"column:name_en;size:255"`
	NameZh string `json:"nameZh" gorm:"column:name_zh;size:255"`
	Value  string `json:"value" gorm:"column:value;type:longtext"`
}

func (About) TableName() string {
	return "gvto_about"
}
