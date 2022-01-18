package dao

import (
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Post struct {
	BaseModel
	ID         string    `gorm:"size:100;not null;primaryKey" json:"id"`
	Title      string    `gorm:"size:200;uniqueIndex;not null" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	Liked      uint      `gorm:"default:0" binding:"-" json:"liked"`
	IsPublic   bool      `gorm:"type:boolean;default:false" binding:"boolean" json:"isPublic"`
	CategoryID uint      `json:"categoryID"`
	Category   *Category `gorm:"foreignkey:CategoryID" binding:"-" json:"category,omitempty"`
	UserID     string    `json:"userID"`
	User       *User     `binding:"-" json:"user,omitempty"`
}

func (m Post) Create() (Post, error) {
	id := uuid.NewV4().String()
	m.ID = id
	if err := db.Create(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (m Post) Save() (Post, error) {
	if err := db.Save(m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (m Post) Update(values interface{}) (Post, error) {
	err := db.Model(&m).Updates(values).Error
	return m, err
}

func UpdatePosts(values interface{}, ids []string) error {
	return db.Model(&Post{}).Where("id IN (?)", ids).Updates(values).Error
}

func DeletePost(id []string) error {
	return db.Delete(&Post{}, id).Error
}

func (m Post) Relations(col string) *gorm.Association {
	return db.Model(&m).Association(col)
}

func FindPost(id string, options map[string]interface{}) (Post, error) {
	var one Post
	if err := db.Scopes(applyQueryOptions(options)).First(&one, "id = ?", id).Error; err != nil {
		return one, err
	}
	return one, nil
}

func FindPosts(options map[string]interface{}) ([]Post, error) {
	var rows []Post
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

func AppExists(id string) (bool, Post) {
	var one Post
	err := db.Where("id = ?", id).First(&one).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func FindAndCountPosts(options map[string]interface{}) ([]Post, int64, error) {
	var rows []Post
	var count int64
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, count, err
	}
	delete(options, "offset")
	delete(options, "limit")
	delete(options, "order")
	delete(options, "join")
	if err := db.Model(&Post{}).Scopes(applyQueryOptions(options)).Count(&count).Error; err != nil {
		return rows, count, err
	}
	return rows, count, nil
}
