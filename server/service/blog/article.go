package blog

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

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

func (s *ArticleService) SearchPublishedBlogs(query string) ([]blogResp.SearchBlogItem, error) {
	query = strings.TrimSpace(query)
	if query == "" || len([]rune(query)) > 20 || hasSearchSpecialChar(query) {
		return nil, errors.New("invalid query")
	}

	likeQuery := "%" + query + "%"
	var blogs []blogModel.Blog
	if err := global.GVA_DB.
		Where("is_published = ? AND (password IS NULL OR password = '')", true).
		Where("title LIKE ? OR content LIKE ?", likeQuery, likeQuery).
		Order("create_time desc").
		Limit(20).
		Find(&blogs).Error; err != nil {
		return nil, err
	}

	results := make([]blogResp.SearchBlogItem, 0, len(blogs))
	for _, item := range blogs {
		results = append(results, blogResp.SearchBlogItem{
			ID:      item.ID,
			Title:   item.Title,
			Content: buildSearchSnippet(item.Content, query),
		})
	}
	return results, nil
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
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		categoryID, err := s.resolveCategoryID(tx, info)
		if err != nil {
			return err
		}
		tagIDs, err := s.resolveTagIDs(tx, info.TagList)
		if err != nil {
			return err
		}

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
			CategoryID:       categoryID,
			IsTop:            info.IsTop,
			Password:         info.Password,
			UserID:           info.UserID,
		}
		if err := tx.Create(&entity).Error; err != nil {
			return err
		}
		return s.syncBlogTags(tx, entity.ID, tagIDs)
	})
}

func (s *AdminArticleService) Update(info blogReq.ArticleUpsert) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		categoryID, err := s.resolveCategoryID(tx, info)
		if err != nil {
			return err
		}
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
			"category_id":        categoryID,
			"is_top":             info.IsTop,
			"password":           info.Password,
			"user_id":            info.UserID,
		}
		if err := tx.Model(&blogModel.Blog{}).Where("id = ?", info.ID).Updates(updates).Error; err != nil {
			return err
		}
		if info.TagList == nil {
			return nil
		}
		tagIDs, err := s.resolveTagIDs(tx, info.TagList)
		if err != nil {
			return err
		}
		return s.syncBlogTags(tx, info.ID, tagIDs)
	})
}

func (s *AdminArticleService) Delete(id uint) error {
	tx := global.GVA_DB.Begin()
	if err := tx.Where("blog_id = ?", id).Delete(&blogModel.BlogTag{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("blog_id = ?", id).Delete(&blogModel.Comment{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&blogModel.Blog{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *AdminArticleService) UpdateTop(id uint, top bool) error {
	return global.GVA_DB.Model(&blogModel.Blog{}).Where("id = ?", id).Update("is_top", top).Error
}

func (s *AdminArticleService) UpdateRecommend(id uint, recommend bool) error {
	return global.GVA_DB.Model(&blogModel.Blog{}).Where("id = ?", id).Update("is_recommend", recommend).Error
}

func (s *AdminArticleService) UpdateVisibility(id uint, info blogReq.BlogVisibility) error {
	updates := map[string]interface{}{
		"update_time": time.Now(),
	}
	if info.Appreciation != nil {
		updates["is_appreciation"] = *info.Appreciation
	}
	if info.Recommend != nil {
		updates["is_recommend"] = *info.Recommend
	}
	if info.CommentEnabled != nil {
		updates["is_comment_enabled"] = *info.CommentEnabled
	}
	if info.Top != nil {
		updates["is_top"] = *info.Top
	}
	if info.Published != nil {
		updates["is_published"] = *info.Published
	}
	if info.Password != nil {
		password := strings.TrimSpace(*info.Password)
		updates["password"] = password
	}
	return global.GVA_DB.Model(&blogModel.Blog{}).Where("id = ?", id).Updates(updates).Error
}

func hasSearchSpecialChar(query string) bool {
	for _, r := range query {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}
		if strings.ContainsRune("-_.,，。！？!?[]()（）", r) {
			continue
		}
		return true
	}
	return false
}

func buildSearchSnippet(content, query string) string {
	lowerContent := strings.ToLower(content)
	lowerQuery := strings.ToLower(query)
	index := strings.Index(lowerContent, lowerQuery)
	if index < 0 {
		runes := []rune(content)
		if len(runes) <= 40 {
			return content
		}
		return string(runes[:40])
	}
	runes := []rune(content)
	start := 0
	matchPos := len([]rune(content[:index]))
	if matchPos > 10 {
		start = matchPos - 10
	}
	end := start + 21
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}

func (s *AdminArticleService) resolveCategoryID(tx *gorm.DB, info blogReq.ArticleUpsert) (uint, error) {
	if info.Cate == nil {
		if info.CategoryID == 0 {
			return 0, errors.New("category is required")
		}
		var category blogModel.Category
		if err := tx.First(&category, info.CategoryID).Error; err != nil {
			return 0, err
		}
		return category.ID, nil
	}

	switch cate := info.Cate.(type) {
	case string:
		name := strings.TrimSpace(cate)
		if name == "" {
			return 0, errors.New("category is required")
		}
		var existing blogModel.Category
		err := tx.Where("category_name = ?", name).First(&existing).Error
		if err == nil {
			return 0, errors.New("category already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
		entity := blogModel.Category{CategoryName: name}
		if err := tx.Create(&entity).Error; err != nil {
			return 0, err
		}
		return entity.ID, nil
	default:
		id, err := mixedValueToUint(cate)
		if err != nil {
			return 0, errors.New("invalid category")
		}
		var category blogModel.Category
		if err := tx.First(&category, id).Error; err != nil {
			return 0, err
		}
		return category.ID, nil
	}
}

func (s *AdminArticleService) resolveTagIDs(tx *gorm.DB, tagList []any) ([]uint, error) {
	if len(tagList) == 0 {
		return nil, nil
	}
	tagIDs := make([]uint, 0, len(tagList))
	seen := make(map[uint]struct{}, len(tagList))
	for _, raw := range tagList {
		switch tag := raw.(type) {
		case string:
			name := strings.TrimSpace(tag)
			if name == "" {
				return nil, errors.New("invalid tag")
			}
			var existing blogModel.Tag
			err := tx.Where("tag_name = ?", name).First(&existing).Error
			if err == nil {
				return nil, errors.New("tag already exists")
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			entity := blogModel.Tag{TagName: name}
			if err := tx.Create(&entity).Error; err != nil {
				return nil, err
			}
			if _, ok := seen[entity.ID]; !ok {
				seen[entity.ID] = struct{}{}
				tagIDs = append(tagIDs, entity.ID)
			}
		default:
			id, err := mixedValueToUint(tag)
			if err != nil {
				return nil, errors.New("invalid tag")
			}
			var existing blogModel.Tag
			if err := tx.First(&existing, id).Error; err != nil {
				return nil, err
			}
			if _, ok := seen[existing.ID]; !ok {
				seen[existing.ID] = struct{}{}
				tagIDs = append(tagIDs, existing.ID)
			}
		}
	}
	return tagIDs, nil
}

func (s *AdminArticleService) syncBlogTags(tx *gorm.DB, blogID uint, tagIDs []uint) error {
	if err := tx.Where("blog_id = ?", blogID).Delete(&blogModel.BlogTag{}).Error; err != nil {
		return err
	}
	if len(tagIDs) == 0 {
		return nil
	}
	relations := make([]blogModel.BlogTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		relations = append(relations, blogModel.BlogTag{
			BlogID: blogID,
			TagID:  tagID,
		})
	}
	return tx.Create(&relations).Error
}

func mixedValueToUint(value any) (uint, error) {
	switch v := value.(type) {
	case uint:
		return v, nil
	case uint8:
		return uint(v), nil
	case uint16:
		return uint(v), nil
	case uint32:
		return uint(v), nil
	case uint64:
		return uint(v), nil
	case int:
		if v < 0 {
			return 0, fmt.Errorf("negative value")
		}
		return uint(v), nil
	case int8:
		if v < 0 {
			return 0, fmt.Errorf("negative value")
		}
		return uint(v), nil
	case int16:
		if v < 0 {
			return 0, fmt.Errorf("negative value")
		}
		return uint(v), nil
	case int32:
		if v < 0 {
			return 0, fmt.Errorf("negative value")
		}
		return uint(v), nil
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("negative value")
		}
		return uint(v), nil
	case float32:
		if v < 0 {
			return 0, fmt.Errorf("negative value")
		}
		return uint(v), nil
	case float64:
		if v < 0 {
			return 0, fmt.Errorf("negative value")
		}
		return uint(v), nil
	case string:
		parsed, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(parsed), nil
	default:
		return 0, fmt.Errorf("unsupported type %T", value)
	}
}
