package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		routerGroup := v1.Group("/users")
		{
			// Public methods
			routerGroup.GET("/:id", handler.readUser)

			// Protected methods: user
			routerGroup.GET("/my", handler.readMyUser)

			// Protected methods: admin or auth-service
			routerGroup.POST("", handler.createUser)
			routerGroup.PUT("/:id", handler.updateUser)
			routerGroup.DELETE("/:id", handler.deleteUser)
		}
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

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and server: %v", err)
		}
	}()

	log.Printf("Wait for interrupt signal")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shutdown server: %v", err)
	}

	log.Printf("Bye bye (˶ᵔᗜᵔ˶)ﾉﾞ")
}
