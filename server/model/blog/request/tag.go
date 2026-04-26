package request

type TagUpsert struct {
	ID      uint    `json:"id"`
	TagName string  `json:"tagName"`
	Color   *string `json:"color"`
}
