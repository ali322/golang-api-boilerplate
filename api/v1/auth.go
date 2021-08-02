package v1

import (
	"app/model"
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
	Repeatpassword string `binding:"required,lt=200" json:"repeat_password"`
	Email          string `binding:"lt=200,email"`
}

func register(c *gin.Context) {
	var request registerRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(err)
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
	c.JSON(http.StatusOK, &gin.H{
		"user": created, "token": token,
	})
}
