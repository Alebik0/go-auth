package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var hmacSecret = os.Getenv("AUTH_SERVICE_HMAC_SECRET")

func NewReverseProxy(target string) *httputil.ReverseProxy {
	remote, err := url.Parse(target)
	if err != nil {
		log.Fatalf("invalid target url: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	// Optional: modify request before sending to backend
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// You can modify headers here
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Host = remote.Host
	}

	return proxy
}

func SetupRouter() (*gin.Engine, error) {
	userServiceHost := os.Getenv("USER_SERVICE_HOST")
	if userServiceHost == "" {
		return nil, fmt.Errorf("USER_SERVICE_HOST is mandatory environment variable")
	}

	userServicePort := os.Getenv("USER_SERVICE_PORT")
	if userServicePort == "" {
		return nil, fmt.Errorf("USER_SERVICE_PORT is mandatory environment variable")
	}

	authServiceHost := os.Getenv("AUTH_SERVICE_HOST")
	if authServiceHost == "" {
		return nil, fmt.Errorf("AUTH_SERVICE_HOST is mandatory environment variable")
	}

	authServicePort := os.Getenv("AUTH_SERVICE_PORT")
	if authServicePort == "" {
		return nil, fmt.Errorf("AUTH_SERVICE_PORT is mandatory environment variable")
	}

	userService := NewReverseProxy(fmt.Sprintf("http://%s:%s", userServiceHost, userServicePort))
	authService := NewReverseProxy(fmt.Sprintf("http://%s:%s", authServiceHost, authServicePort))

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		public := v1.Group("/")
		public.Use(CleanAuthMiddleware())
		{
			auth := public.Group("/auth")
			{
				auth.POST("/register", func(c *gin.Context) {
					authService.ServeHTTP(c.Writer, c.Request)
				})
				auth.POST("/login", func(c *gin.Context) {
					authService.ServeHTTP(c.Writer, c.Request)
				})
				auth.POST("/logout", func(c *gin.Context) {
					authService.ServeHTTP(c.Writer, c.Request)
				})
				auth.POST("/refresh", func(c *gin.Context) {
					authService.ServeHTTP(c.Writer, c.Request)
				})
			}
			users := public.Group("/users")
			{
				users.GET("/my", func(c *gin.Context) {
					userService.ServeHTTP(c.Writer, c.Request)
				})
				users.GET("/:id", func(c *gin.Context) {
					userService.ServeHTTP(c.Writer, c.Request)
				})
			}
		}

		protected := v1.Group("/")
		protected.Use(AuthMiddleware())
		{
			users := public.Group("/users")
			{
				users.POST("", func(c *gin.Context) {
					userService.ServeHTTP(c.Writer, c.Request)
				})
				users.PUT("/:id", func(c *gin.Context) {
					userService.ServeHTTP(c.Writer, c.Request)
				})
				users.DELETE("/:id", func(c *gin.Context) {
					userService.ServeHTTP(c.Writer, c.Request)
				})
			}
		}
	}

	return router, nil
}

func main() {
	log.SetPrefix("[API_GATEWAY] ")

	router, err := SetupRouter()
	if err != nil {
		log.Fatalf("Failed to setup router: %v", err)
	}

	port := os.Getenv("API_GATEWAY_PORT")
	if port == "" {
		port = "8082"
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
