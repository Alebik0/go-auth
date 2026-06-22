package api

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (handler *Handler) getAuthorizationData(context *gin.Context) (jwt.RegisteredClaims, error) {
	const prefix = "Bearer "
	accessToken := context.GetHeader("Authorization")
	if !strings.HasPrefix(accessToken, prefix) {
		return jwt.RegisteredClaims{}, errors.New("missing token")
	}
	accessToken = strings.TrimPrefix(accessToken, prefix)

	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(
		accessToken,
		&claims,
		func(token *jwt.Token) (any, error) {
			return handler.hmacSecret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return jwt.RegisteredClaims{}, fmt.Errorf("invalid token: %v", err)
	}

	return claims, nil
}

func (handler *Handler) isUser(context *gin.Context) bool {
	claims, err := handler.getAuthorizationData(context)
	if err != nil {
		log.Printf("failed get authorization data: %v", err)
	}

	return err == nil && slices.Contains(claims.Audience, "user")
}

func (handler *Handler) isAdmin(context *gin.Context) bool {
	claims, err := handler.getAuthorizationData(context)
	if err != nil {
		log.Printf("failed get authorization data: %v", err)
	}

	return err == nil && slices.Contains(claims.Audience, "admin")
}

func (handler *Handler) getUserID(context *gin.Context) (uint32, error) {
	claims, err := handler.getAuthorizationData(context)
	if err != nil {
		log.Printf("failed get authorization data: %v", err)
		return 0, fmt.Errorf("invalid authorization: %v", err)
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 32)
	if err != nil {
		log.Printf("parser authorization data (subject): %v", err)
		return 0, fmt.Errorf("invalid authorization: %v", err)
	}

	return uint32(userID), nil
}
