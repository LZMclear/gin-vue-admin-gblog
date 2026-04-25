package request

type BlogLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BlogPasswordCheck struct {
	BlogID   uint   `json:"blogId"`
	Password string `json:"password"`
}
