package middleware

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func isMethodAllowed(method string, methods []string) bool {
	for i := 0; i < len(methods); i++ {
		if method == methods[i] {
			return true
		}
	}
	return false
}

func matchRules(rules map[string]string, target string, method string) (bool, error) {
	for rule, m := range rules {
		matched, err := regexp.MatchString(rule, target)
		if err != nil {
			return false, err
		}
		if matched {
			return isMethodAllowed(strings.ToLower(method), strings.Split(strings.ToLower(m), "|")), nil
		}
	}
	return false, nil
}
func decodeToken(tokenStr string, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return token.Claims.(jwt.MapClaims), nil
}

func JWT(unless map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		matched, err := matchRules(unless, c.Request.URL.String(), c.Request.Method)
		if err != nil {
			_ = c.Error(err)
			return
		}
		if !matched {
			_ = c.Error(errors.New("不匹配"))
			return
		}
		headerStr := c.Request.Header.Get("Authorization")
		if headerStr == "" {
			_ = c.Error(errors.New("授权头信息为空"))
			return
		}
		sp := strings.Split(headerStr, "Bearer ")
		if len(sp) < 1 {
			_ = c.Error(errors.New("授权头信息不合法"))
			return
		}
		tokenStr := sp[1]
		env := c.MustGet("env").(map[string]string)
		jwtSecret := env["JWT_SECRET"]
		token, err := decodeToken(tokenStr, jwtSecret)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.Set("user", token["user"])
		c.Next()
	}
}
