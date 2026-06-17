package main

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/alebik0/go-auth/user-service/database"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	application App
}

func (h *Handler) Close() error {
	return h.application.Close()
}

func NewHaldler(app App) Handler {
	return Handler{application: app}
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
func (h *Handler) readMyUser(context *gin.Context) {
	log.Println("Read user")

	log.Println("Load authorized user data")
	userID, err := strconv.ParseInt(context.GetHeader("X-User-ID"), 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "authorization failed"})
		return
	}
	userRoles := strings.Split(context.GetHeader("X-User-Role"), " ")

	log.Println("Check permissions")
	if !slices.Contains(userRoles, "user") {
		context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
		return
	}

	log.Println("Read user")
	userData, err := h.application.Database.ReadUser(uint32(userID))
	if err == database.ErrUserNotFound {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
func (h *Handler) readUser(context *gin.Context) {
	log.Println("Read user")

	log.Println("Load URI parameters")
	var parameters UserRequestParameters
	if err := context.ShouldBindUri(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Read user")
	userData, err := h.application.Database.ReadUser(parameters.ID)
	if err == database.ErrUserNotFound {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
func (h *Handler) updateUser(context *gin.Context) {
	log.Println("Update user")

	log.Println("Load URI parameters")
	var parameters UserRequestParameters
	if err := context.ShouldBindUri(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Load body parameters")
	var body UpdateUserData
	if err := context.ShouldBindBodyWithJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Load authorized user data")
	_, err := strconv.Atoi(context.GetHeader("X-User-ID"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "authorization failed"})
		return
	}
	userRoles := strings.Split(context.GetHeader("X-User-Role"), " ")

	log.Println("Check permissions")
	if !slices.Contains(userRoles, "admin") && !slices.Contains(userRoles, "auth-service") {
		context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
		return
	}

	log.Println("Update user")
	userData, err := h.application.Database.UpdateUser(parameters.ID, body.Name, body.Description)
	if err == database.ErrUserNotFound {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, userData)
}

// @Summary     Create user CRUD operation
// @Description Creates user in the Postgres database and returns created user as a JSON object
// @Tags        User
// @Accept      json
// @Produce     json
// @Param   	parameters body CreateUserData true "Create user parameters"
// @Param 		X-User-ID header string true "Authorized user ID"
// @Param 		X-User-Role header string true "Authorized user roles"
// @Failure     400 {object} APIError "Bad request"
// @Failure     403 {object} APIError "Forbidden"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /api/v1/users [post]
func (h *Handler) createUser(context *gin.Context) {
	log.Println("Create user")

	log.Println("Load body parameters")
	var body CreateUserData
	if err := context.ShouldBindBodyWithJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Load authorized user data")
	_, err := strconv.Atoi(context.GetHeader("X-User-ID"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "authorization failed"})
		return
	}
	userRoles := strings.Split(context.GetHeader("X-User-Role"), " ")

	log.Println("Check permissions")
	if !slices.Contains(userRoles, "admin") && !slices.Contains(userRoles, "auth-service") {
		context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
		return
	}

	log.Println("Create new user")
	userData, err := h.application.Database.CreateUser(body.Name, body.Description)
	if err != nil {
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
func (h *Handler) deleteUser(context *gin.Context) {
	log.Println("Delete user")

	log.Println("Load URI parameters")
	var parameters UserRequestParameters
	if err := context.ShouldBindUri(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Load authorized user data")
	_, err := strconv.Atoi(context.GetHeader("X-User-ID"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "authorization failed"})
		return
	}
	userRoles := strings.Split(context.GetHeader("X-User-Role"), " ")

	log.Println("Check permissions")
	if !slices.Contains(userRoles, "admin") && !slices.Contains(userRoles, "auth-service") {
		context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
		return
	}

	log.Println("Delete user")
	userData, err := h.application.Database.DeleteUser(parameters.ID)
	if err == database.ErrUserNotFound {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, userData)
}
