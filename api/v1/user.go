package v1

import (
	"app/model"
	"app/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type usersQuery struct {
	Key   string `form:"key" binding:"max=10"`
	Page  int    `form:"page,default=1" binding:"min=1" json:"page"`
	Limit int    `form:"limit,default=10" binding:"min=1" json:"limit"`
}

func users(c *gin.Context) {
	var query usersQuery
	if err := c.ShouldBind(&query); err != nil {
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
	if query.Key != "" {
		where = append(where, []interface{}{"username LIKE ?", fmt.Sprintf("%%%s%%", query.Key)})
	}
	users, count, err := model.FindAndCountUsers(map[string]interface{}{
		"where": where,
		// "preload": []string{"Role"},
		"offset": (query.Page - 1) * query.Limit,
		"limit":  query.Limit,
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

type updateUserBody struct {
	Email  string `binding:"omitempty,lt=200,email"`
	Avatar string `binding:"omitempty,url"`
	Memo   string `binding:"omitempty"`
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var body updateUserBody
	if err := c.ShouldBind(&body); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			c.JSON(http.StatusOK, util.Reject(-2, util.TranslateValidatorErrors(errs)))
			return
		} else {
			_ = c.Error(err)
			return
		}
	}
	user := &model.User{Email: body.Email, Avtar: body.Avatar, Memo: body.Memo}
	updated, err := user.Update(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(updated))
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	deleted, err := model.DeleteUser(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(deleted))
}
