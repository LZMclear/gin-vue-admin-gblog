package response

import "time"

type SearchBlogItem struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArchiveBlogItem struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	CreateTime time.Time `json:"createTime"`
	IsTop      bool      `json:"top"`
	Password   string    `json:"password"`
	Privacy    bool      `json:"privacy"`
}
