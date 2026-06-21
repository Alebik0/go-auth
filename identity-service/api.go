package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

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

type UserRequestParameters struct {
	ID uint32 `uri:"id"`
}

type UpdateUserData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateUserData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type APIError struct {
	Error string `json:"error" example:"just a random internal error"`
}

// @Summary     Read current authentificated user
// @Description Reads user ID from authentification data and reads its data from database.
// @Tags        User
// @Accept      json
// @Produce     json
// @Param 		X-User-ID header string true "Authorized user ID"
// @Param 		X-User-Role header string true "Authorized user roles"
// @Failure     400 {object} APIError "Bad request"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /api/v1/users/my [get]
func (handler *Handler) readMyUser(context *gin.Context) {
	log.Printf("Read user")

	log.Printf("Load authorized user data")
	userID, err := strconv.ParseInt(context.GetHeader("X-User-ID"), 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "authorization failed"})
		return
	}
	userRoles := strings.Split(context.GetHeader("X-User-Role"), " ")

	log.Printf("Check permissions")
	if !slices.Contains(userRoles, "user") {
		context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
		return
	}

	log.Printf("Read my user")
	query := "SELECT id, name, description FROM users WHERE id = $1;"
	userData := UserData{}
	err = handler.
		Database.
		QueryRow(query, userID).
		Scan(&userData.ID, &userData.Name, &userData.Description)

	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, userData)
}

// @Summary     Read user CRUD operation
// @Description Reads user from Postgres database and returns as a JSON object
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Failure     400 {object} APIError "Bad request"
// @Failure     404 {object} APIError "User not found"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /api/v1/users/{id} [get]
func (handler *Handler) readUser(context *gin.Context) {
	log.Printf("Read user")

	log.Printf("Load URI parameters")
	var parameters UserRequestParameters
	err := context.ShouldBindUri(&parameters)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Read user by id=%d", parameters.ID)
	query := "SELECT id, name, description FROM users WHERE id = $1;"
	userData := UserData{}
	err = handler.
		Database.
		QueryRow(query, parameters.ID).
		Scan(&userData.ID, &userData.Name, &userData.Description)

	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, userData)
}

// @Summary     Update user CRUD operation
// @Description Updates user in the Postgres database and returns updated user as a JSON object
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Param   	parameters body UpdateUserData true "Update parameters"
// @Param 		X-User-ID header string true "Authorized user ID"
// @Param 		X-User-Role header string true "Authorized user roles"
// @Failure     400 {object} APIError "Bad request"
// @Failure     403 {object} APIError "Forbidden"
// @Failure     404 {object} APIError "User not found"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /api/v1/users/{id} [put]
func (handler *Handler) updateUser(context *gin.Context) {
	log.Printf("Update user")

	log.Printf("Load URI parameters")
	var parameters UserRequestParameters
	if err := context.ShouldBindUri(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Load body parameters")
	var body UpdateUserData
	if err := context.ShouldBindBodyWithJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Load authorized user data")
	_, err := strconv.Atoi(context.GetHeader("X-User-ID"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "authorization failed"})
		return
	}
	userRoles := strings.Split(context.GetHeader("X-User-Role"), " ")

	log.Printf("Check permissions")
	if !slices.Contains(userRoles, "admin") && !slices.Contains(userRoles, "auth-service") {
		context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
		return
	}

	log.Printf("Update user with id=%d", parameters.ID)
	query := "UPDATE users SET name = $1, description = $2 WHERE id=$3 RETURNING id, name, description;"
	userData := UserData{}
	err = handler.
		Database.
		QueryRow(query, body.Name, body.Description, parameters.ID).
		Scan(&userData.ID, &userData.Name, &userData.Description)
	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, userData)
}

// @Summary     Delete user CRUD operation
// @Description Delete user from the Postgres database and returns deleted user as a JSON object
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Param 		X-User-ID header string true "Authorized user ID"
// @Param 		X-User-Role header string true "Authorized user roles"
// @Failure     400 {object} APIError "Bad request"
// @Failure     403 {object} APIError "Forbidden"
// @Failure     404 {object} APIError "User not found"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /api/v1/users/{id} [delete]
func (handler *Handler) deleteUser(context *gin.Context) {
	log.Printf("Delete user")

	log.Printf("Load URI parameters")
	var parameters UserRequestParameters
	if err := context.ShouldBindUri(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Load authorized user data")
	_, err := strconv.Atoi(context.GetHeader("X-User-ID"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "authorization failed"})
		return
	}
	userRoles := strings.Split(context.GetHeader("X-User-Role"), " ")

	log.Printf("Check permissions")
	if !slices.Contains(userRoles, "admin") && !slices.Contains(userRoles, "auth-service") {
		context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
		return
	}

	log.Printf("Delete user")
	query := "DELETE FROM users WHERE id=$1 RETURNING id, name, description;"
	userData := UserData{}
	err = handler.
		Database.
		QueryRow(query, parameters.ID).
		Scan(&userData.ID, &userData.Name, &userData.Description)
	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, userData)
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

	log.Printf("Check if login %s is unique", parameters.Login)
	query := "SELECT id, login, password_hash, profile_id FROM auth WHERE login = $1;"
	authData := AuthData{}
	err := handler.
		Database.
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
	tx, err := handler.Database.Begin()
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
	userData := UserData{
		Name:        parameters.Login,
		Description: "",
	}
	err = handler.
		Database.
		QueryRow(query, parameters.Login, string(passwordHash)).
		Scan(&userData.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Create authentification")
	query = "INSERT INTO auth (login, password_hash, profile_id) VALUES ($1, $2, $3) RETURNING id;"
	authData = AuthData{
		Login:        parameters.Login,
		PasswordHash: string(passwordHash),
		ProfileID:    userData.ID,
	}
	err = handler.
		Database.
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
func (handler *Handler) login(context *gin.Context) {
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
	authData := AuthData{}
	userData := UserData{}
	err := handler.
		Database.
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
func (handler *Handler) refresh(context *gin.Context) {
	// TODO
}
