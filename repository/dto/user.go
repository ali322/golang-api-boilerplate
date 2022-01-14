package dto

import (
	"app/repository/dao"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type QueryUser struct {
	Key       string `form:"key" binding:"max=10"`
	Page      int    `form:"page,default=1" binding:"min=1" json:"page"`
	Limit     int    `form:"limit,default=10" binding:"min=1" json:"limit"`
	SortBy    string `form:"sortBy,default=created_at"`
	SortOrder string `form:"sortOrder,default=desc"`
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
	values := map[string]interface{}{
		"email":      body.Email,
		"avatar":     body.Avatar,
		"memo":       body.Memo,
		"is_actived": body.IsActived,
	}
	values = omitEmpty(values)
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
		Password: body.Password,
	}
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
	if !found.IsActived {
		return found, errors.New("用户未激活")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(body.Password)); err != nil {
		return found, errors.New("密码不正确")
	}
	updated, err := found.Update(map[string]interface{}{"last_logined_at": time.Now()})
	if err != nil {
		return updated, err
	}
	return updated, nil
}

type ChangePassword struct {
	OldPassword    string `binding:"required,lt=100" json:"oldPassword"`
	NewPassword    string `binding:"required,lt=200" json:"newPassword"`
	RepeatPassword string `binding:"required,lt=200" json:"repeatPassword"`
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
	NewPassword    string `binding:"required,lt=200" json:"newPassword"`
	RepeatPassword string `binding:"required,lt=200" json:"repeatPassword"`
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

type ToggleUserActive struct {
	UserID string `binding:"required" json:"userID"`
}

func (body ToggleUserActive) Active() (err error) {
	values := map[string]interface{}{
		"is_actived": true,
	}
	return dao.UpdateUsers(values, strings.Split(body.UserID, ","))
}

func (body ToggleUserActive) Deactive() (err error) {
	values := map[string]interface{}{
		"is_actived": false,
	}
	return dao.UpdateUsers(values, strings.Split(body.UserID, ","))
}
