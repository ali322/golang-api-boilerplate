package v1

import (
	"app/model"
	"app/util"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func generateJWTToken(secret string, user model.User) (string, error) {
	expired := time.Now().Add(time.Hour * 24 * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user, "expire": expired,
	})
	tokenStr, err := token.SignedString([]byte(secret))
	return tokenStr, err
}

type registerBody struct {
	Username       string `binding:"required,lt=100"`
	Password       string `binding:"required,lt=200"`
	Repeatpassword string `binding:"required,lt=200,eqfield=Password" json:"repeat_password"`
	Email          string `binding:"lt=200,email"`
}

func register(c *gin.Context) {
	var body registerBody
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
	exists, _ := model.UserExists(body.Username)
	if exists {
		_ = c.Error(errors.New("用户已存在"))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 4)
	if err != nil {
		_ = c.Error(err)
		return
	}
	user := model.User{
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
	token, err := generateJWTToken(jwtSecret, created)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(map[string]interface{}{
		"user": created, "token": token,
	}))
}

type loginBody struct {
	Username string `binding:"required,lt=100"`
	Password string `binding:"required,lt=200"`
}

func login(c *gin.Context) {
	var body loginBody
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
	exists, found := model.UserExists(body.Username)
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
	token, err := generateJWTToken(jwtSecret, found)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(map[string]interface{}{
		"user": found, "token": token,
	}))
}

type changePasswordBody struct {
	OldPassword    string `binding:"required,lt=100" json:"old_password"`
	NewPassword    string `binding:"required,lt=200" json:"new_password"`
	RepeatPassword string `binding:"required,lt=200" json:"repeat_password"`
}

func changePassword(c *gin.Context) {
	var body changePasswordBody
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
	me := c.GetStringMap("user")
	id := me["id"].(string)
	user, err := model.FindUser(id, nil)
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

type resetPasswordBody struct {
	NewPassword    string `binding:"required,lt=200" json:"new_password"`
	RepeatPassword string `binding:"required,lt=200" json:"repeat_password"`
}

func resetPassword(c *gin.Context) {
	var body resetPasswordBody
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
	user, err := model.FindUser(id, nil)
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
	me := c.GetStringMap("user")
	c.JSON(http.StatusOK, util.Reply(me))
}
