package dao

import (
	"app/util"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	ID            string         `gorm:"size:100;not_null;primary_key" json:"id"`
	Username      string         `gorm:"size:100;unique_index;not_null" json:"username"`
	Password      string         `gorm:"size:200,not_null" json:"-"`
	Email         string         `gorm:"size:200" json:"email"`
	Avatar        string         `gorm:"type:text" json:"avatar"`
	Memo          string         `gorm:"type:text" json:"memo"`
	IsActived     bool           `gorm:"type:boolean;default:true" binding:"-" json:"isActived"`
	LastLoginedAt util.LocalTime `json:"lastLoginedAt"`
}

func (m User) Create() (User, error) {
	id := uuid.NewV4().String()
	m.ID = id
	m.LastLoginedAt = util.LocalTime{Time: time.Now()}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), 4)
	if err != nil {
		return m, err
	}
	m.Password = string(hashedPassword)
	if err := db.Create(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (m User) Save(cols []string) (User, error) {
	tx := db.Session(&gorm.Session{})
	if len(cols) > 0 {
		for _, col := range cols {
			tx = tx.Select(col)
		}
	}
	if err := tx.Updates(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (m User) Update(values interface{}) (User, error) {
	err := db.Model(&m).Updates(values).Error
	return m, err
}

func UpdateUsers(values interface{}, ids []string) error {
	return db.Model(&User{}).Where("id IN (?)", ids).Updates(values).Error
}

func FindByUsername(username string) (bool, User) {
	var one User
	err := db.Where("username = ?", username).First(&one).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func FindByUsernameOrEmail(usernameOrEmail string) (bool, User) {
	var one User
	err := db.Where("username = ?", usernameOrEmail).Or("email = ?", usernameOrEmail).First(&one).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func DeleteUser(id string) (User, error) {
	var one User
	if err := db.Find(&one, "id = ?", id).Error; err != nil {
		return one, err
	}
	err := db.Delete(&one).Error
	return one, err
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
	if err := db.Scopes(applyQueryOptions(options)).First(&one, "id = ?", id).Error; err != nil {
		return one, err
	}
	return one, nil
}

func (user User) Relations(col string) *gorm.Association {
	return db.Model(&user).Association(col)
}
