package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alebik0/go-auth/auth-service/database"
	"github.com/alebik0/go-auth/auth-service/jwt"
	"github.com/gin-gonic/gin"
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

type RefreshResponse struct {
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
	log.Printf("Register new user")

	log.Printf("Parse body data")
	var parameters RegisterRequest
	if err := context.ShouldBindBodyWithJSON(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Check if login is unique")
	_, err := handler.AuthAPI.ReadAuthByLogin(parameters.Login)
	if err == nil {
		context.JSON(http.StatusConflict, gin.H{"error": "Login is already taken"})
		return
	} else if err != database.ErrAuthNotFound {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Create new profile")
	userData, err := handler.UserAPI.CreateUser(parameters.Login, "")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Generate password hash")
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(parameters.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Create new auth")
	_, err = handler.AuthAPI.CreateAuth(parameters.Login, string(passwordHash), userData.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Generate access token")
	accessToken, err := generateAccessToken(handler.hmacSecret, strconv.FormatInt(int64(userData.ID), 10))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Generate refresh token")
	refreshToken, err := generateRefreshToken()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Generate refresh token")
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

	log.Printf("Return refresh and access tokens")
	context.SetCookie(
		"refresh_token",
		refreshToken,
		int((30 * 24 * time.Hour).Seconds()),
		"/api/v1/auth",
		"",
		true,
		true,
	)
	context.IndentedJSON(http.StatusOK, RegisterResponse{AccessToken: accessToken})
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
	log.Printf("Login user")

	log.Printf("Parse body data")
	var parameters RegisterRequest
	if err := context.ShouldBindBodyWithJSON(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Search auth by login %s", parameters.Login)
	authData, err := handler.AuthAPI.ReadAuthByLogin(parameters.Login)
	if err == database.ErrAuthNotFound {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login"})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Check if password is valid %s", parameters.Login)
	err = bcrypt.CompareHashAndPassword([]byte(authData.PasswordHash), []byte(parameters.Password))
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	log.Printf("Search load user by profile id %d", authData.ProfileID)
	userData, err := handler.UserAPI.ReadUser(authData.ProfileID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Generate access token")
	accessToken, err := generateAccessToken(handler.hmacSecret, strconv.FormatInt(int64(userData.ID), 10))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Generate refresh token")
	refreshToken, err := generateRefreshToken()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Save refresh token")
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

	log.Printf("Return access and refresh tokens")
	context.SetCookie(
		"refresh_token",
		refreshToken,
		int((30 * 24 * time.Hour).Seconds()),
		"/api/v1/auth",
		"",
		true,
		true,
	)
	context.IndentedJSON(http.StatusOK, RegisterResponse{AccessToken: accessToken})
}

// @Summary     Logout user
// @Description Logouts current user
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Failure     401 {object} APIError "Unauthorized"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /api/v1/auth/logout [post]
func (handler *Handler) logout(context *gin.Context) {
	log.Printf("Logout user")

	log.Printf("Load refresh_token cookies")
	accessToken, err := context.Cookie("refresh_token")
	if err == http.ErrNoCookie {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "No cookie provided"})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "No cookie provided"})
		return
	}

	log.Printf("Revert refresh token")
	err = handler.JwtAPI.Revert(context.Request.Context(), accessToken)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Remove refresh token from cookies")
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

// @Summary     Refresh access token
// @Description Generate new access token for user
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Failure     401 {object} APIError "Unauthorized"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /api/v1/auth/refresh [post]
func (handler *Handler) refresh(context *gin.Context) {
	log.Printf("Logout user")

	log.Printf("Load refresh_token cookies")
	refreshToken, err := context.Cookie("refresh_token")
	if err == http.ErrNoCookie {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "No cookie provided"})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "No cookie provided"})
		return
	}

	log.Printf("Revert refresh token")
	userID, err := handler.JwtAPI.Get(context.Request.Context(), refreshToken)
	if err == jwt.ErrJwtNotFound {
		context.SetCookie(
			"refresh_token",
			"",
			-1,
			"/api/v1/auth",
			"",
			true,
			true,
		)

		context.JSON(http.StatusUnauthorized, gin.H{"error": "The cookie is out of date"})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := generateAccessToken(handler.hmacSecret, userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Remove refresh token from cookies")
	context.IndentedJSON(http.StatusOK, RefreshResponse{AccessToken: accessToken})
}
