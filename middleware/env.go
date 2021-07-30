package middleware

import "github.com/gin-gonic/gin"

func Env(env map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("env", env)
		c.Next()
	}
}
