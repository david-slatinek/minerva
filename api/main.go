package main

import (
	"context"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"main/config"
	"main/controller"
	"main/database"
	_ "main/docs"
	"main/models"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configFileFlag = flag.String("file", "config.local", "config file name")

//	@title			Song API
//	@version		1.0
//	@description	API for song management
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	David Slatinek
//	@contact.url	https://github.com/david-slatinek

//	@accept		json
//	@produce	json
//	@schemes	http

//	@license.name	Apache-2.0 license
//	@license.url	https://www.apache.org/licenses/LICENSE-2.0

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	flag.Parse()

	cfg, err := config.NewConfig(*configFileFlag)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	db, err := database.NewSong(cfg.ConnectionString)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	songController := controller.NewSong(db)

	gin.SetMode(cfg.Mode)

	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, models.Error{Message: "endpoint not found"})
	})
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, models.Error{Message: "method not allowed"})
	})

	baseGroup := router.Group("api/v1")
	{
		baseGroup.POST("/song", songController.Create)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		log.Println("server is up at: " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("ListenAndServe() error: %s\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Shutdown() error: %s\n", err)
	}

	log.Println("shutting down")
}
