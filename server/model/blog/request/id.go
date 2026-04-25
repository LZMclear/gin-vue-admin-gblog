package request

type IDReq struct {
	ID uint `json:"id" form:"id"`
}

type IDQuery struct {
	ID uint `json:"id" form:"id"`
}
