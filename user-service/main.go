package main

import (
	"log"
	"os"

	"github.com/alebik0/go-auth/user-service/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           User service
// @version         1.0
// @description     User service API, that provides CRUD operations to manage users

// @contact.name Alebik0

// @license.name MIT
// @license.url  https://mit-license.org

// @host     localhost:8080
// @BasePath /api/v1

func SetupRouter(handler Handler) *gin.Engine {
	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		routerGroup := v1.Group("/user")
		routerGroup.POST("", handler.createUser)
		routerGroup.GET("/:id", handler.readUser)
		routerGroup.PUT("/:id", handler.updateUser)
		routerGroup.DELETE("/:id", handler.deleteUser)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}

func main() {
	log.SetPrefix("[USER_SERVICE] ")

	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed create app: %v", err)
	}
	handler := NewHaldler(app)
	defer func() {
		innerErr := handler.Close()
		if innerErr != nil {
			log.Printf("[WARN] Failed to close handler: %v", innerErr)
		}
	}()

	router := SetupRouter(handler)

	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Running server on port %s", port)
	router.Run(":" + port)
}
