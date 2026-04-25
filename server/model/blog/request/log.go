package request

type DateRangePageQuery struct {
	StartDate string `json:"startDate" form:"startDate"`
	EndDate   string `json:"endDate" form:"endDate"`
	Page      int    `json:"page" form:"page"`
	PageSize  int    `json:"pageSize" form:"pageSize"`
}
