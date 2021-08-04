package v1

import (
	"app/model"
	"app/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type usersRequest struct {
	Key   string `form:"key" binding:"max=10"`
	Page  int    `form:"page,default=1" binding:"min=1" json:"page"`
	Limit int    `form:"limit,default=10" binding:"min=1" json:"limit"`
}

func users(c *gin.Context) {
	var request usersRequest
	if err := c.ShouldBind(&request); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			c.JSON(http.StatusOK, util.Reject(-2, util.TranslateValidatorErrors(errs)))
			return
		} else {
			_ = c.Error(err)
			return
		}
	}
	where := make([][]interface{}, 0)
	if request.Key != "" {
		where = append(where, []interface{}{"username LIKE ?", fmt.Sprintf("%%%s%%", request.Key)})
	}
	users, count, err := model.FindAndCountUsers(map[string]interface{}{
		"where": where,
		// "preload": []string{"Role"},
		"offset": (request.Page - 1) * request.Limit,
		"limit":  request.Limit,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(map[string]interface{}{
		"count": count,
		"rows":  users,
	}))
}

func user(c *gin.Context) {
	id := c.Param("id")
	user, err := model.FindUser(id, map[string]interface{}{})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(user))
}
