package model

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID            string    `gorm:"size:100;not_null;primary_key" json:"id"`
	Username      string    `gorm:"size:100;unique_index;not_null"`
	Password      string    `gorm:"size:200,not_null"`
	Email         string    `gorm:"size:200"`
	Avtar         string    `gorm:"type:text"`
	Memo          string    `gorm:"type:text"`
	LastLoginedAt time.Time `time_format:"2016-01-02 15:04:05" json:"last_logined_at"`
	CreatedAt     time.Time `binding:"-"`
	UpdatedAt     time.Time `binding:"-"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	id := uuid.NewV4().String()
	_ = scope.SetColumn("id", id)
	_ = scope.SetColumn("lastLoginedAt", time.Now())
	return nil
}

func (user User) Create() (User, error) {
	if err := db.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
