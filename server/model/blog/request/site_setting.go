package request

type SiteSettingUpsert struct {
	ID     uint    `json:"id"`
	NameEn *string `json:"nameEn"`
	NameZh *string `json:"nameZh"`
	Value  *string `json:"value"`
	Type   *int    `json:"type"`
}

type SiteSettingBatchUpdate struct {
	Settings  []SiteSettingUpsert `json:"settings"`
	DeleteIDs []uint              `json:"deleteIds"`
}
