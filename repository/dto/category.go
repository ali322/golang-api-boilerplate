package dto

import (
	"app/repository/dao"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type NewCategory struct {
	Name        string `binding:"omitempty,lt=200" json:"name"`
	Description string `json:"description"`
	ParentID    *int64 `binding:"omitempty,numeric,gt=0" json:"parentID"`
}

func (body *NewCategory) Create(userID string) (dao.Category, error) {
	m := dao.Category{
		Name: body.Name, Description: body.Description,
	}
	var parentID uint = 1
	if body.ParentID != nil {
		parentID = uint(*body.ParentID)
	}
	parent, err := dao.FindCategory(parentID, nil)
	if err != nil {
		return m, err
	}
	return m.Create(&parent)
}

type UpdateCategory struct {
	Name        string `binding:"omitempty,lt=200" json:"name"`
	Description string `json:"description"`
}

func (body *UpdateCategory) Save(id uint) (dao.Category, error) {
	m, err := dao.FindCategory(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m, errors.New("分类不存在")
		} else {
			return m, err
		}
	}
	values := map[string]interface{}{
		"name":        body.Name,
		"description": body.Description,
	}
	values = omitEmpty(values)
	return m.Update(values)
}

type QueryCategory struct {
	Key       string `form:"key" binding:"max=10" json:"key"`
	Page      int    `form:"page,default=1" binding:"min=1" json:"page"`
	Limit     int    `form:"limit,default=10" binding:"min=1" json:"limit"`
	SortBy    string `form:"sortBy,default=created_at" binding:"oneof=created_at updated_at"`
	SortOrder string `form:"sortOrder,default=desc" binding:"oneof=asc desc"`
}

func (query *QueryCategory) Find() ([]dao.Category, int64, error) {
	where := make([][]interface{}, 0)
	if query.Key != "" {
		where = append(where, []interface{}{"name LIKE ?", fmt.Sprintf("%%%s%%", query.Key)})
	}
	return dao.FindAndCountCategories(map[string]interface{}{
		"where":   where,
		"preload": []string{"Posts"},
		"offset":  (query.Page - 1) * query.Limit,
		"limit":   query.Limit,
		"order":   fmt.Sprintf("%s %s", query.SortBy, query.SortOrder),
	})
}

func isPostExists(row dao.Post, rows []dao.Post) bool {
	for i := 0; i < len(rows); i++ {
		if row.ID == rows[i].ID {
			return true
		}
	}
	return false
}

type IOCategory struct {
	CategoryID uint   `binding:"required" json:"categoryID"`
	PostID     string `binding:"required" json:"postID"`
}

func (body *IOCategory) In() (dao.Category, error) {
	m, err := dao.FindCategory(body.CategoryID, map[string]interface{}{
		"preload": []string{"Posts"},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m, errors.New("分类不存在")
		} else {
			return m, err
		}
	}
	rows, err := dao.FindPosts(map[string]interface{}{
		"where": strings.Split(body.PostID, ","),
	})
	if err != nil {
		return m, err
	}
	next := make([]dao.Post, 0)
	for _, row := range rows {
		if !isPostExists(row, m.Posts) {
			next = append(next, row)
		}
	}
	err = m.Add(next)
	if err != nil {
		return m, err
	}
	return m, nil
}

func (body *IOCategory) Out() (dao.Category, error) {
	m, err := dao.FindCategory(body.CategoryID, map[string]interface{}{
		"preload": []string{"Posts"},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m, errors.New("分类不存在")
		} else {
			return m, err
		}
	}
	rows, err := dao.FindPosts(map[string]interface{}{
		"where": strings.Split(body.PostID, ","),
	})
	if err != nil {
		return m, err
	}
	next := make([]dao.Post, 0)
	left := make([]dao.Post, 0)
	for _, post := range m.Posts {
		if isPostExists(post, rows) {
			next = append(next, post)
		} else {
			left = append(left, post)
		}
	}
	err = m.Remove(next)
	if err != nil {
		return m, err
	}
	m.Posts = left
	return m, nil
}

type MovePost struct {
	From   uint   `binding:"required" json:"from"`
	To     uint   `binding:"required" json:"to"`
	PostID string `binding:"required" json:"postID"`
}

func (body MovePost) Move() (dao.Category, error) {
	from, err := dao.FindCategory(body.From, map[string]interface{}{
		"preload": []string{"Posts"},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return from, errors.New("来源分类不存在")
		} else {
			return from, err
		}
	}
	to, err := dao.FindCategory(body.To, map[string]interface{}{
		"preload": []string{"Posts"},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return to, errors.New("目标分类不存在")
		} else {
			return to, err
		}
	}
	rows, err := dao.FindPosts(map[string]interface{}{
		"where": strings.Split(body.PostID, ","),
	})
	if err != nil {
		return to, err
	}
	fromRows := make([]dao.Post, 0)
	for _, post := range from.Posts {
		if isPostExists(post, rows) {
			fromRows = append(fromRows, post)
		}
	}
	toRows := make([]dao.Post, 0)
	for _, row := range rows {
		if !isPostExists(row, to.Posts) {
			toRows = append(toRows, row)
		}
	}
	err = dao.MoveCategory(&from, &to, fromRows, toRows)
	if err != nil {
		return to, err
	}
	return to, nil
}

type DeleteCategory struct {
	ID string `binding:"omitempty" json:"id"`
}

func (body DeleteCategory) Delete() error {
	ids := make([]uint, 0)
	for _, id := range strings.Split(body.ID, ",") {
		id, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		ids = append(ids, uint(id))
	}
	return dao.DeleteCategories(ids)
}

type MoveCategory struct {
	ID string `binding:"omitempty" json:"id"`
}

func (body *MoveCategory) Move(parent *dao.Category) (err error) {
	rows, err := dao.FindCategories(map[string]interface{}{
		"where": strings.Split(body.ID, ","),
	})
	for _, row := range rows {
		err := row.MoveTo(parent)
		if err != nil {
			return err
		}
	}
	return
}
