package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		errs := c.Errors
		if len(errs) == 0 {
			return
		}
		err := errs.Last()
		// errs, ok := err.(validator.ValidationErrors)
		c.JSON(http.StatusOK, &gin.H{
			"code": -1, "message": err.Error(),
		})
		c.Abort()
	}
}
