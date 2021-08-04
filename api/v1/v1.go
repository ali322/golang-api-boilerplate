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
		v1.POST("register", register)
		v1.POST("login", login)
		v1.POST("change/:id/password", changePassword)
		v1.POST("reset/:id/password", resetPassword)
		v1.GET("user", users)
		v1.GET("user/:id", user)
		v1.PUT("user/:id", updateUser)
		v1.DELETE("user/:id", deleteUser)
	}
}
