package blog

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
	blogResp "github.com/flipped-aurora/gin-vue-admin/server/model/blog/response"
	"gorm.io/gorm"
)

type ArticleService struct{}
type AdminArticleService struct{}

func (s *ArticleService) GetPublishedList(info blogReq.ArticleSearch) (list []blogModel.Blog, total int64, err error) {
	db := global.GVA_DB.Model(&blogModel.Blog{}).Where("is_published = ?", true)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	offset := info.PageSize * (info.Page - 1)
	err = db.Order("is_top desc, create_time desc").Limit(info.PageSize).Offset(offset).Find(&list).Error
	return
}

func (s *ArticleService) GetPublishedByID(id uint) (blog blogModel.Blog, err error) {
	err = global.GVA_DB.Where("id = ? AND is_published = ?", id, true).First(&blog).Error
	return
}

func (s *ArticleService) GetPublishedByIDWithToken(id uint, rawToken string) (blog blogModel.Blog, err error) {
	err = global.GVA_DB.Where("id = ? AND is_published = ?", id, true).First(&blog).Error
	if err != nil {
		return
	}
	if _, err = ensureBlogReadable(blog, rawToken); err != nil {
		return blog, err
	}
	blog.Password = nil
	_ = global.GVA_DB.Model(&blogModel.Blog{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + 1")).Error
	blog.Views++
	return blog, nil
}

func (s *AdminArticleService) GetList(info blogReq.AdminArticleSearch) (list []blogModel.Blog, total int64, err error) {
	db := global.GVA_DB.Model(&blogModel.Blog{})
	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}
	if info.CategoryID != nil && *info.CategoryID > 0 {
		db = db.Where("category_id = ?", *info.CategoryID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	offset := info.PageSize * (info.Page - 1)
	err = db.Order("create_time desc").Limit(info.PageSize).Offset(offset).Find(&list).Error
	return
}

func (s *AdminArticleService) GetByID(id uint) (blog blogModel.Blog, err error) {
	err = global.GVA_DB.First(&blog, id).Error
	return
}

func (s *AdminArticleService) GetCategoryAndTag() (res blogResp.CategoryAndTagResponse, err error) {
	err = global.GVA_DB.Order("id asc").Find(&res.Categories).Error
	if err != nil {
		return
	}
	err = global.GVA_DB.Order("id asc").Find(&res.Tags).Error
	return
}

func (s *AdminArticleService) Create(info blogReq.ArticleUpsert) error {
	now := time.Now()
	entity := blogModel.Blog{
		Title:            info.Title,
		FirstPicture:     info.FirstPicture,
		Content:          info.Content,
		Description:      info.Description,
		IsPublished:      info.IsPublished,
		IsRecommend:      info.IsRecommend,
		IsAppreciation:   info.IsAppreciation,
		IsCommentEnabled: info.IsCommentEnabled,
		CreateTime:       now,
		UpdateTime:       now,
		Views:            info.Views,
		Words:            info.Words,
		ReadTime:         info.ReadTime,
		CategoryID:       info.CategoryID,
		IsTop:            info.IsTop,
		Password:         info.Password,
		UserID:           info.UserID,
	}
	return global.GVA_DB.Create(&entity).Error
}

func (s *AdminArticleService) Update(info blogReq.ArticleUpsert) error {
	updates := map[string]interface{}{
		"title":              info.Title,
		"first_picture":      info.FirstPicture,
		"content":            info.Content,
		"description":        info.Description,
		"is_published":       info.IsPublished,
		"is_recommend":       info.IsRecommend,
		"is_appreciation":    info.IsAppreciation,
		"is_comment_enabled": info.IsCommentEnabled,
		"update_time":        time.Now(),
		"views":              info.Views,
		"words":              info.Words,
		"read_time":          info.ReadTime,
		"category_id":        info.CategoryID,
		"is_top":             info.IsTop,
		"password":           info.Password,
		"user_id":            info.UserID,
	}
	return global.GVA_DB.Model(&blogModel.Blog{}).Where("id = ?", info.ID).Updates(updates).Error
}

func (s *AdminArticleService) Delete(id uint) error {
	tx := global.GVA_DB.Begin()
	if err := tx.Where("blog_id = ?", id).Delete(&blogModel.BlogTag{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&blogModel.Blog{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
