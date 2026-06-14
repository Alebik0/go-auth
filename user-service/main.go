package main

import (
	"log"
	"os"

	"github.com/alebik0/go-auth/user-service/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var app App

// @title           User service
// @version         1.0
// @description     User service API, that provides CRUD operations to manage users

// @contact.name Alebik0

// @license.name MIT
// @license.url  https://mit-license.org

// @host     localhost:8080
// @BasePath /api/v1

func main() {
	log.SetPrefix("[USER_SERVICE] ")

	var err error
	app, err = NewApp()
	if err != nil {
		log.Fatalf("Failed create app: %v", err)
	}
	defer app.Close()

	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		routerGroup := v1.Group("/user")
		routerGroup.POST("", createUser)
		routerGroup.GET("/:id", readUser)
		routerGroup.PUT("/:id", updateUser)
		routerGroup.DELETE("/:id", deleteUser)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Running server on port %s", port)
	router.Run(":" + port)
}
