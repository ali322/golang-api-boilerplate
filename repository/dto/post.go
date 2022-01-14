package dto

import (
	"app/repository/dao"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type NewPost struct {
	Title      string `binding:"omitempty,lt=200" json:"title"`
	Content    string `json:"content"`
	CategoryID uint   `binding:"required,numeric,gt=0" json:"categoryID"`
	IsPublic   *bool  `binding:"omitempty" json:"isPublic"`
}

func (body *NewPost) Create(userID string) (dao.Post, error) {
	m := dao.Post{
		Title: body.Title, Content: body.Content,
		CategoryID: body.CategoryID,
	}
	if body.IsPublic != nil {
		m.IsPublic = *body.IsPublic
	}
	created, err := m.Create()
	if err != nil {
		return created, err
	}
	return created, nil
}

type UpdatePost struct {
	Title      string `binding:"omitempty,lt=200" json:"title"`
	Content    string `json:"content"`
	CategoryID *uint  `binding:"omitempty,numeric,gt=0" json:"categoryID"`
	IsPublic   bool   `binding:"omitempty" json:"isPublic"`
}

func (body *UpdatePost) Save(id string) (dao.Post, error) {
	m, err := dao.FindPost(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m, errors.New("文章不存在")
		} else {
			return m, err
		}
	}
	if body.CategoryID != nil {
		exists, _ := dao.CategoryExists(*body.CategoryID)
		if !exists {
			return m, errors.New("分类不存在")
		}
	}
	values := map[string]interface{}{
		"title":       body.Title,
		"content":     body.Content,
		"category_id": body.CategoryID,
		"is_public":   body.IsPublic,
	}
	values = omitEmpty(values)
	updated, err := m.Update(values)
	if err != nil {
		return updated, err
	}
	return updated, nil
}

type QueryPost struct {
	Key       string `form:"key" binding:"max=10"`
	Page      int    `form:"page,default=1" binding:"min=1" json:"page"`
	Limit     int    `form:"limit,default=10" binding:"min=1" json:"limit"`
	SortBy    string `form:"sortBy,default=created_at"`
	SortOrder string `form:"sortOrder,default=desc"`
}

func (query *QueryPost) Find() ([]dao.Post, int64, error) {
	where := make([][]interface{}, 0)
	if query.Key != "" {
		where = append(where, []interface{}{"title LIKE ?", fmt.Sprintf("%%%s%%", query.Key)})
	}
	return dao.FindAndCountPosts(map[string]interface{}{
		"where": where,
		// "preload": []string{"Role"},
		"offset": (query.Page - 1) * query.Limit,
		"limit":  query.Limit,
		"order":  fmt.Sprintf("%s %s", query.SortBy, query.SortOrder),
	})
}
