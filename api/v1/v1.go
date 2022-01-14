package v1

import (
	"app/middleware"
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
		v1.POST("public/register", register)
		v1.POST("public/login", login)
		v1.POST("change/password", changePassword)
		v1.POST("reset/:id/password", resetPassword)
		v1.GET("me", me)
		v1.GET("user", users)
		v1.GET("user/:id", user)
		v1.PUT("user/:id", updateUser)
		v1.DELETE("user/:id", deleteUser)
		v1.POST("active/user", activeUser)
		v1.DELETE("active/user", deactiveUser)

		v1.GET("public/connect/message", ConnectWebsocket)
		v1.GET("public/disconnect/message", DisconnectWebsocket)

		v1.POST("post", createPost)
		v1.PUT("post/:id", updatePost)
		v1.GET("public/post/:id", post)
		v1.GET("public/post", middleware.Cache(), posts)

		v1.POST("category", createCategory)
		v1.PUT("category/:id", updateCategory)
		v1.GET("public/category/:id", category)
		v1.GET("public/category", categories)
		v1.DELETE("category", deleteCategory)
		v1.POST("category/to/:id", moveCategory)
		v1.POST("category/post", addToCategory)
		v1.DELETE("category/post", removeFromCategory)
		v1.PUT("category/post", movePost)
	}
}
