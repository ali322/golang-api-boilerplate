package model

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            string    `gorm:"size:100;not_null;primary_key" json:"id"`
	Username      string    `gorm:"size:100;unique_index;not_null"`
	Password      string    `gorm:"size:200,not_null"`
	Email         string    `gorm:"size:200"`
	Avtar         string    `gorm:"type:text"`
	Memo          string    `gorm:"type:text"`
	LastLoginedAt time.Time `time_format:"2016-01-02 15:04:05" json:"last_logined_at"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	id := uuid.NewV4().String()
	tx.Statement.SetColumn("id", id)
	tx.Statement.SetColumn("lastLoginedAt", time.Now())
	return nil
}

func (user User) Create() (User, error) {
	if err := db.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (user User) Save() error {
	return db.Save(&user).Error
}

func (user User) Update(fields map[string]interface{}) error {
	return db.Model(&user).Updates(fields).Error
}

func UserExists(username string) (bool, User) {
	var one User
	err := db.Where("username = ?", username).First(&one).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func DeleteUser(id string) error {
	var one User
	if err := db.Find(&one, id).Error; err != nil {
		return err
	}
	return db.Delete(&one).Error
	// return db.Model(&one).Update(&User{IsDeleted: true}).Error
}

func FindAndCountUsers(options map[string]interface{}) ([]User, int64, error) {
	var rows []User
	var count int64
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, count, err
	}
	delete(options, "offset")
	delete(options, "limit")
	if err := db.Model(&User{}).Scopes(applyQueryOptions(options)).Count(&count).Error; err != nil {
		return rows, count, err
	}
	return rows, count, nil
}

func FindUsers(options map[string]interface{}) ([]User, error) {
	var rows []User
	if err := db.Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

func FindUser(id string, options map[string]interface{}) (User, error) {
	var one User
	if err := db.Scopes(applyQueryOptions(options)).Find(&one, "id = ?", id).Error; err != nil {
		return one, err
	}
	return one, nil
}

func (user User) Relations(col string) *gorm.Association {
	return db.Model(&user).Association(col)
}
