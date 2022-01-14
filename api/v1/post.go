package v1

import (
	"app/repository/dao"
	"app/repository/dto"
	"app/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createPost(c *gin.Context) {
	var body dto.NewPost
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	auth := c.GetStringMap("auth")
	id := auth["id"].(string)
	created, err := body.Create(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(created))
}

func updatePost(c *gin.Context) {
	id := c.Param("id")
	var body dto.UpdatePost
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	saved, err := body.Save(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(saved))
}

func post(c *gin.Context) {
	id := c.Param("id")
	found, err := dao.FindPost(id, map[string]interface{}{
		"preload": []string{"Category"},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(found))
}

func posts(c *gin.Context) {
	var query dto.QueryPost
	if err := c.ShouldBind(&query); err != nil {
		_ = c.Error(err)
		return
	}
	rows, count, err := query.Find()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(map[string]interface{}{
		"count": count,
		"rows":  rows,
	}))
}
