package blog

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
)

type SiteService struct{}

func (s *SiteService) GetSiteInfo() (map[string]interface{}, error) {
	result := map[string]interface{}{}

	var siteSettings []blogModel.SiteSetting
	if err := global.GVA_DB.Order("id asc").Find(&siteSettings).Error; err != nil {
		return nil, err
	}

	var categories []blogModel.Category
	if err := global.GVA_DB.Order("id asc").Find(&categories).Error; err != nil {
		return nil, err
	}

	var tags []blogModel.Tag
	if err := global.GVA_DB.Order("id asc").Find(&tags).Error; err != nil {
		return nil, err
	}

	var newBlogs []blogModel.Blog
	if err := global.GVA_DB.Where("is_published = ?", true).Order("create_time desc").Limit(5).Find(&newBlogs).Error; err != nil {
		return nil, err
	}

	var randomBlogs []blogModel.Blog
	if err := global.GVA_DB.Where("is_published = ? AND is_recommend = ?", true, true).Order("rand()").Limit(5).Find(&randomBlogs).Error; err != nil {
		return nil, err
	}

	result["siteSettings"] = siteSettings
	result["categoryList"] = categories
	result["tagList"] = tags
	result["newBlogList"] = newBlogs
	result["randomBlogList"] = randomBlogs

	return result, nil
}
