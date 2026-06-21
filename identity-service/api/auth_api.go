package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alebik0/go-auth/identity-service/data"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789"
)

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
func (handler *Handler) Register(context *gin.Context) {
	log.Printf("Register new user")

	log.Printf("Parse body data")
	var parameters RegisterRequest
	if err := context.ShouldBindBodyWithJSON(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Check if login %s is unique", parameters.Login)
	query := "SELECT id, login, password_hash, profile_id FROM auth WHERE login = $1;"
	authData := data.AuthData{}
	err := handler.
		database.
		QueryRow(query, parameters.Login).
		Scan(&authData.ID, &authData.Login, &authData.PasswordHash, &authData.ProfileID)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Login %s is free", parameters.Login)
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		context.JSON(http.StatusConflict, gin.H{"error": "Login is already taken"})
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

	log.Printf("Create new profile")

	log.Printf("Begin transaction")
	tx, err := handler.database.Begin()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		err = tx.Rollback()
		if err != nil {
			log.Printf("[WARN] Failed to rollback transaction: %v", err)
		}
	}()

	log.Printf("Create profile")
	query = "INSERT INTO users (name, description) VALUES ($1, $2) RETURNING id;"
	userData := data.UserData{
		Name:        parameters.Login,
		Description: "",
	}
	err = handler.
		database.
		QueryRow(query, parameters.Login, string(passwordHash)).
		Scan(&userData.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Create authentification")
	query = "INSERT INTO auth (login, password_hash, profile_id) VALUES ($1, $2, $3) RETURNING id;"
	authData = data.AuthData{
		Login:        parameters.Login,
		PasswordHash: string(passwordHash),
		ProfileID:    userData.ID,
	}
	err = handler.
		database.
		QueryRow(query, authData.Login, authData.PasswordHash, authData.ProfileID).
		Scan(&userData.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tx.Commit()
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
	err = handler.jwtAPI.Save(
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
		"/api/v1",
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
func (handler *Handler) Login(context *gin.Context) {
	log.Printf("Login user")

	log.Printf("Parse body data")
	var parameters RegisterRequest
	if err := context.ShouldBindBodyWithJSON(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Search auth by login %s", parameters.Login)
	query := `
	SELECT
		a.id AS auth_id,
		a.login,
		a.password_hash,
		a.profile_id,
		u.id AS user_id,
		u.name,
		u.description
	FROM auth a
	JOIN users u ON a.profile_id = u.id
	WHERE a.login = $1;
	`
	authData := data.AuthData{}
	userData := data.UserData{}
	err := handler.
		database.
		QueryRow(query, parameters.Login).
		Scan(
			&authData.ID,
			&authData.Login,
			&authData.PasswordHash,
			&authData.ProfileID,
			&userData.ID,
			&userData.Name,
			&userData.Description,
		)
	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
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
	err = handler.jwtAPI.Save(
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
		"/api/v1",
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
func (handler *Handler) Logout(context *gin.Context) {
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
	err = handler.jwtAPI.Revert(context.Request.Context(), accessToken)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Remove refresh token from cookies")
	context.SetCookie(
		"refresh_token",
		"",
		-1,
		"/api/v1",
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
func (handler *Handler) Refresh(context *gin.Context) {
	// TODO
}
