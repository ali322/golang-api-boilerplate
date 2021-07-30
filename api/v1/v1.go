package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, &gin.H{
				"code": 0, "message": "pong",
			})
		})
	}
}
