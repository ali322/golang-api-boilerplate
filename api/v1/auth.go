package v1

import (
	"app/repository/dao"
	"app/repository/dto"
	"app/util"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	created, err := body.Create()
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
	found, err := body.Login()
	if err != nil {
		_ = c.Error(err)
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
	updated, err := body.ChangePassword(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(updated))
}

func resetPassword(c *gin.Context) {
	var body dto.ResetPassword
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
	updated, err := body.ResetPassword(id)
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
