package request

type PageQuery struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type ToggleReq struct {
	ID    uint `json:"id" form:"id"`
	Value bool `json:"value" form:"value"`
}

type NicknameReq struct {
	Nickname string `json:"nickname" form:"nickname"`
}

type ContentReq struct {
	Content string `json:"content"`
}

type IDUriReq struct {
	ID uint `uri:"id" binding:"required"`
}
