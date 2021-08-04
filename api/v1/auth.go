package v1

import (
	"app/model"
	"app/util"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func generateJWTToken(secret string, user model.User) (string, error) {
	expired := time.Now().Add(time.Hour * 24 * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user, "expire": expired,
	})
	tokenStr, err := token.SignedString([]byte(secret))
	return tokenStr, err
}

type registerRequest struct {
	Username       string `binding:"required,lt=100"`
	Password       string `binding:"required,lt=200"`
	Repeatpassword string `binding:"required,lt=200,eqfield=Password" json:"repeat_password"`
	Email          string `binding:"lt=200,email"`
}

func register(c *gin.Context) {
	var request registerRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(err)
		return
	}
	exists, _ := model.UserExists(request.Username)
	if exists {
		_ = c.Error(errors.New("用户已存在"))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 4)
	if err != nil {
		_ = c.Error(err)
		return
	}
	user := model.User{
		Username: request.Username,
		Password: string(hashedPassword),
		Email:    request.Email,
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

type loginRequest struct {
	Username string `binding:"required,lt=100"`
	Password string `binding:"required,lt=200"`
}

func login(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(err)
		return
	}
	exists, found := model.UserExists(request.Username)
	if !exists {
		_ = c.Error(errors.New("用户不存在"))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(request.Password)); err != nil {
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

type changePasswordRequest struct {
	OldPassword    string `binding:"required;lt=200"`
	NewPassword    string `binding:"required;lt=200"`
	RepeatPassword string `binding:"required;lt=200"`
}

func changePassword(c *gin.Context) {
	var request changePasswordRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(err)
		return
	}
	id := c.Param("id")
	user, err := model.FindUser(id, nil)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)); err != nil {
		_ = c.Error(err)
		return
	}
	if request.NewPassword != request.RepeatPassword {
		_ = c.Error(errors.New("重复密码不匹配"))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 4)
	if err != nil {
		_ = c.Error(err)
		return
	}
	user.Password = string(hashedPassword)
	if err := user.Save(); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(user))
}

type resetPasswordRequest struct {
	NewPassword    string `binding:"required;lt=200"`
	RepeatPassword string `binding:"required;lt=200"`
}

func resetPassword(c *gin.Context) {
	var request resetPasswordRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(err)
		return
	}
	id := c.Param("id")
	user, err := model.FindUser(id, nil)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if request.NewPassword != request.RepeatPassword {
		_ = c.Error(errors.New("重复密码不匹配"))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 4)
	if err != nil {
		_ = c.Error(err)
		return
	}
	user.Password = string(hashedPassword)
	if err := user.Save(); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(user))
}
