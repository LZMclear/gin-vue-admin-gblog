package blog

type SiteSetting struct {
	ID     uint    `json:"id" gorm:"column:id;primaryKey"`
	NameEn *string `json:"nameEn,omitempty" gorm:"column:name_en;size:255"`
	NameZh *string `json:"nameZh,omitempty" gorm:"column:name_zh;size:255"`
	Value  *string `json:"value,omitempty" gorm:"column:value;type:longtext"`
	Type   *int    `json:"type,omitempty" gorm:"column:type"`
}

func (SiteSetting) TableName() string {
	return "gvto_site_setting"
}
