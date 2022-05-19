package main

import (
	"app/api"
	"app/lib/config"
	"app/lib/logger"
	"app/lib/ws"
	"app/middleware"
	"app/repository/dao"
	"app/util"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func setupApp() *gin.Engine {
	logger.Init(config.App.LogDir)
	defer logger.Logger.Sync()
	app := gin.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())
	app.Use(middleware.Error())
	app.Use(middleware.Cors())
	app.Use(middleware.JWT(map[string]string{
		"public": "post|get",
	}))
	util.InitTranslator(config.App.Locale)
	util.RegisterValidatorTranslations(config.App.Locale)
	go ws.WebsocketServer.Start()
	dao.Init(config.Database.URL)
	api.ApplyRoutes(app)
	return app
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config.Read()
	app := setupApp()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.App.Port),
		Handler: app,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen: %s\n", err)
		}
	}()
	// err := app.Run(fmt.Sprintf(":%s", config.App.Port))
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	<-ctx.Done()
	stop()
	log.Println("shutdown gracefully, press ctrl+c force shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dao.Close(); err != nil {
		log.Fatal("failed to close db: ", err)
	}
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("failed to shutdown server: ", err)
	}
	log.Println("server exiting")
	// shutdown.NewHook().Close(
	// 	func() {
	// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 		defer cancel()

	// 		if err := server.Shutdown(ctx); err != nil {
	// 			log.Fatal("failed to shutdown server: ", err)
	// 		}
	// 	},
	// 	func() {
	// 		if err := dao.Close(); err != nil {
	// 			log.Fatal("failed to close db: ", err)
	// 		}
	// 	},
	// )
}
