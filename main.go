package main

import (
	"app/api"
	"app/middleware"
	"app/repository/dao"
	"app/util"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func pwd() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return path, nil
}

func setupApp(env map[string]string) *gin.Engine {
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(middleware.Recovery())
	app.Use(middleware.Error())
	app.Use(middleware.Env(env))
	app.Use(middleware.JWT(map[string]string{
		// "auth": "post|get",
		"login":    "post",
		"register": "post",
	}))
	util.InitTranslator(env["LOCALE"])
	util.RegisterValidatorTranslations(env["LOCALE"])
	dao.Init(env)
	api.ApplyRoutes(app)
	return app
}

func main() {
	var path string = ""
	if os.Getenv("GIN_MODE") == "release" {
		path, _ = pwd()
	} else {
		gin.SetMode(gin.DebugMode)
	}
	env, err := godotenv.Read(filepath.Join(path, ".env"))
	if err != nil {
		log.Fatal("failed to read .env")
	}
	app := setupApp(env)
	err = app.Run(fmt.Sprintf(":%s", env["APP_PORT"]))
	if err != nil {
		log.Fatal(err.Error())
	}
}
