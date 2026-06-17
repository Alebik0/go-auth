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
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		const prefix = "Bearer "
		accessToken := context.GetHeader("Authorization")
		if !strings.HasPrefix(accessToken, prefix) {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		accessToken = strings.TrimPrefix(accessToken, prefix)

		var claims jwt.RegisteredClaims
		_, err := jwt.ParseWithClaims(
			accessToken,
			&claims,
			func(token *jwt.Token) (any, error) {
				return hmacSecret, nil
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("invalid token: %v", err)})
			return
		}

		context.Request.Header.Set("X-User-ID", claims.Subject)
		context.Request.Header.Set("X-User-Role", strings.Join(claims.Audience, " "))
		context.Next()
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	userService := NewReverseProxy("http://localhost:8080")
	authService := NewReverseProxy("http://localhost:8081")

	v1 := router.Group("/api/v1")
	{
		public := v1.Group("/")
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

	return router
}

func main() {
	log.SetPrefix("[API_GATEWAY] ")

	router := SetupRouter()

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
