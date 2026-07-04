package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alebik0/go-auth/identity-service/api"
	"github.com/alebik0/go-auth/identity-service/docs"
	"github.com/alebik0/go-auth/identity-service/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Identity service
// @version         1.0
// @description     Identity service API, that provides authentification functionality, session control via JWT tokens and user management.

// @license.name MIT
// @license.url  https://mit-license.org

// @host     localhost:8081
// @BasePath /api/v1

func SetupRouter(handler api.Handler) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true, // if using cookies
		MaxAge: 12 * time.Hour,
	}))

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("login", handler.Login)
			auth.POST("logout", handler.Logout)
			auth.POST("register", handler.Register)
			auth.POST("refresh", handler.Refresh)
		}
		users := v1.Group("/users")
		{
			users.GET("/:id", handler.ReadUser)

			protected := users.Group("")
			protected.Use(middleware.RequireAuth(handler.HmacSecret))
			protected.Use(middleware.RequireAnyRole(middleware.UserRole))
			{
				protected.GET("my", handler.ReadMyUser)
				protected.PUT("/:id", handler.UpdateUser)
				protected.DELETE("/:id", handler.DeleteUser)
			}
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}

func main() {
	log.SetPrefix("[IDENTITY_SERVICE] ")

	database, err := NewDatabaseDependency()
	if err != nil {
		log.Fatalf("Failed create database dependency: %v", err)
	}

	cache, err := NewCacheDependency()
	if err != nil {
		log.Fatalf("Failed create jwt dependency: %v", err)
	}

	hmacSecret := os.Getenv("IDENTITY_SERVICE_HMAC_SECRET")
	if hmacSecret == "" {
		log.Fatalf("IDENTITY_SERVICE_HMAC_SECRET is mandatory environment variable")
	}

	handler, err := api.NewHandler(
		database,
		cache,
		[]byte(hmacSecret),
	)
	if err != nil {
		log.Fatalf("Failed create app: %v", err)
	}
	defer func() {
		innerErr := handler.Close()
		if innerErr != nil {
			log.Printf("[WARN] Failed to close handler: %v", innerErr)
		}
	}()

	router := SetupRouter(handler)

	port := os.Getenv("IDENTITY_SERVICE_PORT")
	if port == "" {
		port = "8081"
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
