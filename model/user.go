package model

import (
	"api-boilerplate/dao"
)

type User struct {
	ID       string `gorm:"size:100;not_null;primary_key" json:"id"`
	Username string `gorm:"size:100;unique_index;not_null" binding:"required,lt=50" json:"username"`
}

func (user User) Create() (User, error) {
	if err := dao.Database.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
