package util

import (
	"github.com/gin-gonic/gin"
)

func Reply(data interface{}) *gin.H {
	return &gin.H{
		"code": 0, "data": data,
	}
}

func Reject(code int, data interface{}) *gin.H {
	return &gin.H{
		"code": code, "message": data,
	}
}
