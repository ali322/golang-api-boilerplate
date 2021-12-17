package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				print(err)
				c.JSON(http.StatusOK, &gin.H{
					"code": -1, "message": err,
				})
			}
		}()
		c.Next()
	}
}
