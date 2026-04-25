package request

type CommentSearch struct {
	Page     int   `json:"page" form:"page"`
	BlogID   *uint `json:"blogId" form:"blogId"`
	PageNum  int   `json:"pageNum" form:"pageNum"`
	PageSize int   `json:"pageSize" form:"pageSize"`
}

type CommentAdminSearch struct {
	Page     *int  `json:"page" form:"page"`
	BlogID   *uint `json:"blogId" form:"blogId"`
	PageNum  int   `json:"pageNum" form:"pageNum"`
	PageSize int   `json:"pageSize" form:"pageSize"`
}

type CommentCreate struct {
	Nickname        string  `json:"nickname"`
	Email           string  `json:"email"`
	Content         string  `json:"content"`
	Avatar          string  `json:"avatar"`
	IP              *string `json:"ip"`
	IsPublished     bool    `json:"isPublished"`
	IsAdminComment  bool    `json:"isAdminComment"`
	Page            int     `json:"page"`
	IsNotice        bool    `json:"isNotice"`
	BlogID          *uint   `json:"blogId"`
	ParentCommentID int64   `json:"parentCommentId"`
	Website         *string `json:"website"`
	QQ              *string `json:"qq"`
}

type CommentUpdate struct {
	ID              uint    `json:"id"`
	Nickname        string  `json:"nickname"`
	Email           string  `json:"email"`
	Content         string  `json:"content"`
	Avatar          string  `json:"avatar"`
	IP              *string `json:"ip"`
	IsPublished     bool    `json:"isPublished"`
	IsAdminComment  bool    `json:"isAdminComment"`
	Page            int     `json:"page"`
	IsNotice        bool    `json:"isNotice"`
	BlogID          *uint   `json:"blogId"`
	ParentCommentID int64   `json:"parentCommentId"`
	Website         *string `json:"website"`
	QQ              *string `json:"qq"`
}
