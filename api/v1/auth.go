package v1

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func generateJWTToken(secret string, user wswdxszdao.User) (string, error) {
	expired := time.Now().Add(time.Hour * 24 * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user, "expire": expired,
	})
	tokenStr, err := token.SignedString([]byte(secret))
	return tokenStr, err
}
