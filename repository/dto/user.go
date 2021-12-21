package dto

import (
	"app/repository/dao"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type QueryUser struct {
	Key       string `form:"key" binding:"max=10"`
	Page      int    `form:"page,default=1" binding:"min=1" json:"page"`
	Limit     int    `form:"limit,default=10" binding:"min=1" json:"limit"`
	SortBy    string `form:"sort_by,default=created_at"`
	SortOrder string `form:"sort_order,default=desc"`
}

func (query *QueryUser) Find() ([]dao.User, int64, error) {
	where := make([][]interface{}, 0)
	if query.Key != "" {
		where = append(where, []interface{}{"username LIKE ?", fmt.Sprintf("%%%s%%", query.Key)})
	}
	return dao.FindAndCountUsers(map[string]interface{}{
		"where": where,
		// "preload": []string{"Role"},
		"offset": (query.Page - 1) * query.Limit,
		"limit":  query.Limit,
		"order":  fmt.Sprintf("%s %s", query.SortBy, query.SortOrder),
	})
}

type UpdateUser struct {
	Email     string `binding:"omitempty,lt=200,email"`
	Avatar    string `binding:"omitempty,url"`
	Memo      string `binding:"omitempty"`
	IsActived bool   `binding:"omitempty" json:"is_actived"`
}

func (body *UpdateUser) Save(id string) (dao.User, error) {
	user, err := dao.FindUser(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("用户不存在")
		} else {
			return user, err
		}
	}
	values := make(map[string]interface{})
	if body.Email != "" {
		values["email"] = body.Email
	}
	if body.Avatar != "" {
		values["avatar"] = body.Avatar
	}
	if body.Memo != "" {
		values["memo"] = body.Memo
	}
	println("%v", body.IsActived)
	values["is_actived"] = body.IsActived
	return user.Update(values)
}

type RegisterUser struct {
	Username       string `binding:"required,lt=100"`
	Password       string `binding:"required,lt=200"`
	Repeatpassword string `binding:"required,lt=200,eqfield=Password" json:"repeat_password"`
	Email          string `binding:"lt=200,email"`
}

func (body *RegisterUser) Create() (dao.User, error) {
	user := dao.User{
		Username: body.Username,
		Email:    body.Email,
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 4)
	if err != nil {
		return user, err
	}
	user.Password = string(hashedPassword)
	return user.Create()
}

type LoginUser struct {
	UsernameOrEmail string `binding:"required,lt=100" json:"username_or_email"`
	Password        string `binding:"required,lt=200"`
}

func (body *LoginUser) Login() (dao.User, error) {
	exists, found := dao.FindByUsernameOrEmail(body.UsernameOrEmail)
	if !exists {
		return found, errors.New("用户不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(body.Password)); err != nil {
		return found, errors.New("密码不正确")
	}
	return found.Update(map[string]interface{}{"last_logined_at": time.Now()})
}

type ChangePassword struct {
	OldPassword    string `binding:"required,lt=100" json:"old_password"`
	NewPassword    string `binding:"required,lt=200" json:"new_password"`
	RepeatPassword string `binding:"required,lt=200" json:"repeat_password"`
}

func (body *ChangePassword) ChangePassword(id string) (dao.User, error) {
	user, err := dao.FindUser(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("用户不存在")
		} else {
			return user, err
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword)); err != nil {
		return user, errors.New("旧密码不正确")
	}
	if body.NewPassword != body.RepeatPassword {
		return user, errors.New("重复密码不匹配")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 4)
	if err != nil {
		return user, err
	}
	return user.Update(map[string]interface{}{"password": string(hashedPassword)})
}

type ResetPassword struct {
	NewPassword    string `binding:"required,lt=200" json:"new_password"`
	RepeatPassword string `binding:"required,lt=200" json:"repeat_password"`
}

func (body *ResetPassword) ResetPassword(id string) (dao.User, error) {
	user, err := dao.FindUser(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("用户不存在")
		} else {
			return user, err
		}
	}
	if body.NewPassword != body.RepeatPassword {
		return user, errors.New("重复密码不匹配")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 4)
	if err != nil {
		return user, err
	}
	return user.Update(map[string]interface{}{"password": string(hashedPassword)})
}
