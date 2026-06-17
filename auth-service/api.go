package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alebik0/go-auth/auth-service/database"
	userservice "github.com/alebik0/go-auth/auth-service/user-service"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789"
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type APIError struct {
	Error string `json:"error" example:"just a random internal error"`
}

// @Summary     Registers new user
// @Description Registers new user with the provided login and password, login must be unique
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param   	parameters body RegisterRequest true "Register parameters"
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

	accessToken, err := generateAccessToken(handler.hmacSecret, userData)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = handler.JwtAPI.Save(
		context.Request.Context(),
		refreshToken,
		strconv.FormatInt(int64(userData.ID), 10),
		time.Now().Add(30*24*time.Hour),
	)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
	context.IndentedJSON(http.StatusOK, registerResponse)
}

func generateAccessToken(hmacSecret []byte, userData userservice.UserData) (string, error) {
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
	token := make([]byte, 32)

	for i := range token {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		token[i] = alphabet[n.Int64()]
	}

	return string(token), nil
}

// @Summary     Login user
// @Description Logins user with the provided login and password and returns, checks if provided password is correct via checking password hash match
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param   	parameters body LoginRequest true "Login parameters"
// @Failure     400 {object} APIError "Bad request"
// @Failure     401 {object} APIError "Unauthorized"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /api/v1/auth/login [post]
func (handler *Handler) login(context *gin.Context) {
	log.Println("Login user")

	var parameters RegisterRequest
	if err := context.ShouldBindBodyWithJSON(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authData, err := handler.AuthAPI.ReadAuthByLogin(parameters.Login)
	if err == database.ErrAuthNotFound {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login"})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(authData.PasswordHash), []byte(parameters.Password))
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	userData, err := handler.UserAPI.ReadUser(authData.ProfileID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := generateAccessToken(handler.hmacSecret, userData)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = handler.JwtAPI.Save(
		context.Request.Context(),
		refreshToken,
		strconv.FormatInt(int64(userData.ID), 10),
		time.Now().Add(30*24*time.Hour),
	)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
	context.IndentedJSON(http.StatusOK, registerResponse)
}

// @Summary     Logout user
// @Description Logouts current user
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Failure     400 {object} APIError "Bad request"
// @Failure     401 {object} APIError "Unauthorized"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /api/v1/auth/logout [post]
func (handler *Handler) logout(context *gin.Context) {
	log.Println("Logout user")

	const prefix = "Bearer "
	accessToken := context.GetHeader("Authorization")
	if !strings.HasPrefix(accessToken, prefix) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token: should start from Bearer"})
		return
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
		context.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid JWT token: parse failure: %v", err)})
		return
	}

	if time.Now().After(claims.ExpiresAt.Time) || time.Now().Before(claims.NotBefore.Time) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token: expired"})
		return
	}

	userIDSubject := claims.Subject
	// roles := claims.Audience

	userID, err := strconv.Atoi(userIDSubject)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token: invalid Subject"})
		return
	}

	log.Printf("User is logined: userID=%d", userID)

	// TODO: remove refreshToken

	context.SetCookie(
		"refresh_token",
		"",
		-1,
		"/api/v1/auth",
		"",
		true,
		true,
	)
	context.IndentedJSON(http.StatusOK, "")
}
