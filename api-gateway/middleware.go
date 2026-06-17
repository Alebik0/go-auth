package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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

func CleanAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Request.Header.Set("X-User-ID", "")
		context.Request.Header.Set("X-User-Role", "")
		context.Next()
	}
}
