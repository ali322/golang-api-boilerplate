package v1

import (
	"app/repository/dao"
	"app/repository/dto"
	"app/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createCategory(c *gin.Context) {
	var body dto.NewCategory
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

func updateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	var body dto.UpdateCategory
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	saved, err := body.Save(uint(id))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(saved))
}

func deleteCategory(c *gin.Context) {
	var body dto.DeleteCategory
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	err := body.Delete()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func category(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	found, err := dao.FindCategoryHierarchy(uint(id), map[string]interface{}{
		"preload": []string{"Posts"},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(found))
}

func categories(c *gin.Context) {
	var query dto.QueryCategory
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

func moveCategory(c *gin.Context) {
	var body dto.MoveCategory
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	parentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	parent, err := dao.FindCategory(uint(parentID), nil)
	if err != nil {
		_ = c.Error(err)
		return
	}
	err = body.Move(&parent)
	// err = folder.MoveTo(&parent)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func addToCategory(c *gin.Context) {
	var body dto.IOCategory
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	folder, err := body.In()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(folder))
}

func removeFromCategory(c *gin.Context) {
	var body dto.IOCategory
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	folder, err := body.Out()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(folder))
}

func movePost(c *gin.Context) {
	var body dto.MovePost
	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	folder, err := body.Move()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, util.Reply(folder))
}
