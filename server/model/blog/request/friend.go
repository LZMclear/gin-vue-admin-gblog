package request

type FriendUpsert struct {
	ID          uint   `json:"id"`
	Nickname    string `json:"nickname"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Avatar      string `json:"avatar"`
	IsPublished bool   `json:"isPublished"`
	Views       int    `json:"views"`
}
