package response

import blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"

type Option struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

type CategoryAndTagResponse struct {
	Categories []blogModel.Category `json:"categories"`
	Tags       []blogModel.Tag      `json:"tags"`
}
