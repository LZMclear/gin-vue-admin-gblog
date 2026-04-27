package request

type ArticleSearch struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type SearchBlogQuery struct {
	Query string `json:"query" form:"query"`
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
	Cate             any     `json:"cate"`
	TagList          []any   `json:"tagList"`
}

type BlogVisibility struct {
	Appreciation   *bool   `json:"appreciation"`
	Recommend      *bool   `json:"recommend"`
	CommentEnabled *bool   `json:"commentEnabled"`
	Top            *bool   `json:"top"`
	Published      *bool   `json:"published"`
	Password       *string `json:"password"`
}
