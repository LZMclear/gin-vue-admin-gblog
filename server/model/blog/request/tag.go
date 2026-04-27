package request

type TagUpsert struct {
	ID      uint    `json:"id"`
	TagName string  `json:"tagName"`
	Color   *string `json:"color"`
}

type TagBlogSearch struct {
	TagName  string `json:"tagName" form:"tagName"`
	Page     int    `json:"page" form:"page"`
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}
