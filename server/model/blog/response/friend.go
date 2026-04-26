package response

type FriendInfoResponse struct {
	Content        string `json:"content"`
	CommentEnabled bool   `json:"commentEnabled"`
}
