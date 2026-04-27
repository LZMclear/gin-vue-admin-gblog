package request

type CategoryUpsert struct {
	ID           uint   `json:"id"`
	CategoryName string `json:"categoryName"`
}

type CategoryBlogSearch struct {
	CategoryName string `json:"categoryName" form:"categoryName"`
	Page         int    `json:"page" form:"page"`
	PageNum      int    `json:"pageNum" form:"pageNum"`
	PageSize     int    `json:"pageSize" form:"pageSize"`
}
