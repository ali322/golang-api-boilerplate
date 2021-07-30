package main

import (
	"api-boilerplate/api"
	"api-boilerplate/middleware"
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
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(middleware.Recovery())
	app.Use(middleware.Error())
	app.Use(middleware.Env(env))
	// app.Use(middleware.JWT(map[string]string{
	// 	"auth": "post|get",
	// }))
	api.ApplyRoutes(app)
	port := env["APP_PORT"]
	err = app.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err.Error())
	}
}
