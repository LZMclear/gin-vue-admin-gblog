package request

type ArticleSearch struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type AdminArticleSearch struct {
	Title      string `json:"title" form:"title"`
	CategoryID *uint  `json:"categoryId" form:"categoryId"`
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"pageSize" form:"pageSize"`
}

type ArticleUpsert struct {
	ID               uint    `json:"id"`
	Title            string  `json:"title"`
	FirstPicture     string  `json:"firstPicture"`
	Content          string  `json:"content"`
	Description      string  `json:"description"`
	IsPublished      bool    `json:"isPublished"`
	IsRecommend      bool    `json:"isRecommend"`
	IsAppreciation   bool    `json:"isAppreciation"`
	IsCommentEnabled bool    `json:"isCommentEnabled"`
	Views            int     `json:"views"`
	Words            int     `json:"words"`
	ReadTime         int     `json:"readTime"`
	CategoryID       uint    `json:"categoryId"`
	IsTop            bool    `json:"isTop"`
	Password         *string `json:"password"`
	UserID           *uint   `json:"userId"`
}
