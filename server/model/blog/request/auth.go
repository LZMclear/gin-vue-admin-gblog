package request

type BlogPasswordCheck struct {
	BlogID   uint   `json:"blogId"`
	Password string `json:"password"`
}
