package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alebik0/go-auth/auth-service/database"
	userservice "github.com/alebik0/go-auth/auth-service/user-service"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	hmacSecret = []byte(os.Getenv("AUTH_SERVICE_HMAC_SECRET"))
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
}

type Claims struct {
	UserID           uint32   `json:"user_id"`
	Roles            []string `json:"roles"`
	RegisteredClaims jwt.RegisteredClaims
}

type APIError struct {
	Error string `json:"error" example:"just a random internal error"`
}

var refreshTokens map[uint32]string = make(map[uint32]string)

// @Summary     Registers new user
// @Description Registers new user with the provided login and password and returns, login must be unique
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param   	parameters body RegisterRequest true "Login parameters"
// @Failure     400 {object} APIError "Bad request"
// @Failure     409 {object} APIError "Conflict: login is already taken"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /api/v1/auth/register [post]
func (handler *Handler) register(context *gin.Context) {
	log.Println("Register new user")

	var parameters RegisterRequest
	if err := context.ShouldBindBodyWithJSON(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := handler.AuthAPI.ReadAuthByLogin(parameters.Login)
	if err == nil {
		context.JSON(http.StatusConflict, gin.H{"error": "Login is already taken"})
		return
	} else if err != database.ErrAuthNotFound {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userData, err := handler.UserAPI.CreateUser(parameters.Login, "")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(parameters.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = handler.AuthAPI.CreateAuth(parameters.Login, string(passwordHash), userData.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := generateAccessToken(userData)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// FIXME: rewrite using Redis
	refreshTokens[userData.ID] = refreshToken
	// err = h.refreshRepo.Save(
	// 	r.Context(),
	// 	user.ID,
	// 	refreshToken,
	// 	time.Now().Add(30*24*time.Hour),
	// )
	// if err != nil {
	// 	http.Error(w, "internal error", http.StatusInternalServerError)
	// 	return
	// }

	registerResponse := RegisterResponse{
		AccessToken: accessToken,
	}

	context.SetCookie(
		"refresh_token",
		refreshToken,
		int((30 * 24 * time.Hour).Seconds()),
		"/api/v1/auth",
		"",
		true,
		true,
	)
	context.Header("Content-Type", "application/json")
	context.IndentedJSON(http.StatusOK, registerResponse)
}

func generateAccessToken(userData userservice.UserData) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:   "auth-service",
			Subject:  strconv.FormatInt(int64(userData.ID), 10),
			Audience: []string{"USER"},
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(15 * time.Minute),
			),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "auth_" + strconv.FormatInt(int64(userData.ID), 10),
		},
	)

	return token.SignedString(hmacSecret)
}

func generateRefreshToken() (string, error) {
	secret := make([]byte, 32) // 256-bit key

	_, err := rand.Read(secret)
	if err != nil {
		return "", fmt.Errorf("failed generate random string: %w", err)
	}

	return string(secret), nil
}
