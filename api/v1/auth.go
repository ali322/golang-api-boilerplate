package v1

import (
	"app/repository/dao"
	"app/repository/dto"
	"app/util"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func register(c *gin.Context) {
	var body dto.RegisterUser
	if err := c.ShouldBind(&body); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			c.JSON(http.StatusOK, util.Reject(-2, util.TranslateValidatorErrors(errs)))
			return
		} else {
			_ = c.Error(err)
			return
		}
	}
	exists, _ := dao.FindByUsername(body.Username)
	if exists {
		_ = c.Error(errors.New("用户已存在"))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 4)
	if err != nil {
		_ = c.Error(err)
		return
	}
	user := dao.User{
		Username: body.Username,
		Password: string(hashedPassword),
		Email:    body.Email,
	}
	created, err := user.Create()
	if err != nil {
		_ = c.Error(err)
		return
	}
	env := c.MustGet("env").(map[string]string)
	jwtSecret := env["JWT_SECRET"]
	token, err := util.GenerateToken(jwtSecret, map[string]interface{}{
		"id": created.ID, "username": created.Username,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(map[string]interface{}{
		"user": created, "token": token,
	}))
}

func login(c *gin.Context) {
	var body dto.LoginUser
	if err := c.ShouldBind(&body); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			c.JSON(http.StatusOK, util.Reject(-2, util.TranslateValidatorErrors(errs)))
			return
		} else {
			_ = c.Error(err)
			return
		}
	}
	exists, found := dao.FindByUsernameOrEmail(body.UsernameOrEmail)
	if !exists {
		_ = c.Error(errors.New("用户不存在"))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(body.Password)); err != nil {
		_ = c.Error(errors.New("密码不正确"))
		return
	}
	env := c.MustGet("env").(map[string]string)
	jwtSecret := env["JWT_SECRET"]
	token, err := util.GenerateToken(jwtSecret, map[string]interface{}{
		"id": found.ID, "username": found.Username,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(map[string]interface{}{
		"user": found, "token": token,
	}))
}

func changePassword(c *gin.Context) {
	var body dto.ChangePassword
	if err := c.ShouldBind(&body); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			c.JSON(http.StatusOK, util.Reject(-2, util.TranslateValidatorErrors(errs)))
			return
		} else {
			_ = c.Error(err)
			return
		}
	}
	auth := c.GetStringMap("auth")
	id := auth["id"].(string)
	user, err := dao.FindUser(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(errors.New("用户不存在"))
			return
		} else {
			_ = c.Error(err)
			return
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword)); err != nil {
		_ = c.Error(errors.New("旧密码不正确"))
		return
	}
	if body.NewPassword != body.RepeatPassword {
		_ = c.Error(errors.New("重复密码不匹配"))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 4)
	if err != nil {
		_ = c.Error(err)
		return
	}
	user.Password = string(hashedPassword)
	updated, err := user.Update(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(updated))
}

func resetPassword(c *gin.Context) {
	var body dto.ChangePassword
	if err := c.ShouldBind(&body); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			c.JSON(http.StatusOK, util.Reject(-2, util.TranslateValidatorErrors(errs)))
			return
		} else {
			_ = c.Error(err)
			return
		}
	}
	id := c.Param("id")
	user, err := dao.FindUser(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(errors.New("用户不存在"))
			return
		} else {
			_ = c.Error(err)
			return
		}
	}
	if body.NewPassword != body.RepeatPassword {
		_ = c.Error(errors.New("重复密码不匹配"))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 4)
	if err != nil {
		_ = c.Error(err)
		return
	}
	user.Password = string(hashedPassword)
	updated, err := user.Update(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(updated))
}

func me(c *gin.Context) {
	auth := c.GetStringMap("auth")
	c.JSON(http.StatusOK, util.Reply(auth))
}
