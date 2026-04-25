package blog

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	blogModel "github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	blogReq "github.com/flipped-aurora/gin-vue-admin/server/model/blog/request"
	"gorm.io/gorm"
)

type CommentService struct{}

func (s *CommentService) GetPublicList(info blogReq.CommentSearch) (result map[string]interface{}, err error) {
	db := global.GVA_DB.Model(&blogModel.Comment{}).Where("page = ?", info.Page)
	if info.BlogID != nil {
		db = db.Where("blog_id = ?", *info.BlogID)
	}

	var allCount int64
	if err = db.Count(&allCount).Error; err != nil {
		return nil, err
	}

	var openCount int64
	if err = db.Where("is_published = ?", true).Count(&openCount).Error; err != nil {
		return nil, err
	}

	if info.PageNum <= 0 {
		info.PageNum = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	offset := info.PageSize * (info.PageNum - 1)
	var list []blogModel.Comment
	if err = db.Order("create_time desc").Limit(info.PageSize).Offset(offset).Find(&list).Error; err != nil {
		return nil, err
	}

	result = map[string]interface{}{
		"allComment":   allCount,
		"closeComment": allCount - openCount,
		"comments": map[string]interface{}{
			"list":     list,
			"page":     info.PageNum,
			"pageSize": info.PageSize,
			"total":    allCount,
		},
	}
	return
}

func (s *CommentService) Create(info blogReq.CommentCreate) error {
	now := time.Now()
	entity := blogModel.Comment{
		Nickname:        info.Nickname,
		Email:           info.Email,
		Content:         info.Content,
		Avatar:          info.Avatar,
		CreateTime:      &now,
		IP:              info.IP,
		IsPublished:     info.IsPublished,
		IsAdminComment:  info.IsAdminComment,
		Page:            info.Page,
		IsNotice:        info.IsNotice,
		BlogID:          info.BlogID,
		ParentCommentID: info.ParentCommentID,
		Website:         info.Website,
		QQ:              info.QQ,
	}
	return global.GVA_DB.Create(&entity).Error
}

func (s *CommentService) GetAdminList(info blogReq.CommentAdminSearch) (list []blogModel.Comment, total int64, err error) {
	db := global.GVA_DB.Model(&blogModel.Comment{})
	if info.Page != nil {
		db = db.Where("page = ?", *info.Page)
	}
	if info.BlogID != nil {
		db = db.Where("blog_id = ?", *info.BlogID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.PageNum <= 0 {
		info.PageNum = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	offset := info.PageSize * (info.PageNum - 1)
	err = db.Order("create_time desc").Limit(info.PageSize).Offset(offset).Find(&list).Error
	return
}

func (s *CommentService) UpdatePublished(id uint, published bool) error {
	return global.GVA_DB.Model(&blogModel.Comment{}).Where("id = ?", id).Update("is_published", published).Error
}

func (s *CommentService) UpdateNotice(id uint, notice bool) error {
	return global.GVA_DB.Model(&blogModel.Comment{}).Where("id = ?", id).Update("is_notice", notice).Error
}

func (s *CommentService) Delete(id uint) error {
	return global.GVA_DB.Where("id = ? OR parent_comment_id = ?", id, id).Delete(&blogModel.Comment{}).Error
}

func (s *CommentService) Update(info blogReq.CommentUpdate) error {
	return global.GVA_DB.Model(&blogModel.Comment{}).Where("id = ?", info.ID).Updates(map[string]interface{}{
		"nickname":          info.Nickname,
		"email":             info.Email,
		"content":           info.Content,
		"avatar":            info.Avatar,
		"ip":                info.IP,
		"is_published":      info.IsPublished,
		"is_admin_comment":  info.IsAdminComment,
		"page":              info.Page,
		"is_notice":         info.IsNotice,
		"blog_id":           info.BlogID,
		"parent_comment_id": info.ParentCommentID,
		"website":           info.Website,
		"qq":                info.QQ,
	}).Error
}

func (s *CommentService) GetBlogIDAndTitle() ([]blogModel.Blog, error) {
	var list []blogModel.Blog
	err := global.GVA_DB.Select("id,title").Order("create_time desc").Find(&list).Error
	return list, err
}

func (s *CommentService) GetCommentsByBlogID(id uint) ([]blogModel.Comment, error) {
	var list []blogModel.Comment
	err := global.GVA_DB.Where("blog_id = ?", id).Order("create_time desc").Find(&list).Error
	return list, err
}

func (s *CommentService) Count() (int64, error) {
	var total int64
	err := global.GVA_DB.Model(&blogModel.Comment{}).Count(&total).Error
	return total, err
}

func (s *CommentService) CountToday() (int64, error) {
	var total int64
	err := global.GVA_DB.Model(&blogModel.Comment{}).Where("date(create_time) = curdate()").Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return total, nil
}
