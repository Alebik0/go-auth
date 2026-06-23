package middleware

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthorizedUserKey   string = "AuthorizedUserKey"
	AuthorizationHeader string = "Authorization"
	AuthorizationPrefix string = "Bearer "
)

type Role string

type AuthorizedUser struct {
	UserID uint32
	Roles  []Role
}

func RequireAuth(hmacSecret []byte) gin.HandlerFunc {
	return func(context *gin.Context) {
		log.Printf("Middleware: check authorization")

		accessToken := context.GetHeader(AuthorizationHeader)
		if !strings.HasPrefix(accessToken, AuthorizationPrefix) {
			log.Printf("[ERROR] not authorized: bearer prefix expected, actual: %s", accessToken)
			context.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
			return
		}

		accessToken = strings.TrimPrefix(accessToken, AuthorizationPrefix)

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
			log.Printf("[ERROR] not authorized: %v", err)
			context.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
			return
		}

		userId, err := strconv.ParseUint(claims.Subject, 10, 32)
		if err != nil {
			log.Printf("[ERROR] not authorized: %v", err)
			context.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
			return
		}

		roles := make([]Role, len(claims.Audience))
		for index, element := range claims.Audience {
			roles[index] = Role(element)
		}

		authorizedUser := AuthorizedUser{
			UserID: uint32(userId),
			Roles:  roles,
		}
		context.Set(AuthorizedUserKey, authorizedUser)
		context.Next()
	}
}

func RequireAnyRole(roles ...Role) gin.HandlerFunc {
	return func(context *gin.Context) {
		log.Printf("Middleware: check any permission")

		authorizedUser := context.MustGet(AuthorizedUserKey).(AuthorizedUser)
		for _, role := range roles {
			if slices.Contains(authorizedUser.Roles, role) {
				context.Next()
				return
			}
		}

		context.JSON(http.StatusForbidden, gin.H{"error": "not enough permissions"})
	}
}
