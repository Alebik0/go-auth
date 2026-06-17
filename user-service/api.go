package main

import (
	"log"
	"net/http"

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

// @Summary     Read user CRUD operation
// @Description Reads user from Postgres database and returns as a JSON object
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Failure     400 {object} APIError "Bad request"
// @Failure     404 {object} APIError "User not found"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /user/{id} [get]
func (h *Handler) readUser(context *gin.Context) {
	log.Println("Read user")

	var parameters UserRequestParameters
	if err := context.ShouldBindUri(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := h.application.Database.ReadUser(parameters.ID)
	if err == database.UserNotFound {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		context.IndentedJSON(http.StatusOK, userData)
	}
}

// @Summary     Update user CRUD operation
// @Description Updates user in the Postgres database and returns updated user as a JSON object
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Param   	parameters body UpdateUserData true "Update parameters"
// @Failure     400 {object} APIError "Bad request"
// @Failure     404 {object} APIError "User not found"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /user/{id} [put]
func (h *Handler) updateUser(context *gin.Context) {
	log.Println("Update user")

	var parameters UserRequestParameters
	if err := context.ShouldBindUri(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body UpdateUserData
	if err := context.ShouldBindBodyWithJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := h.application.Database.UpdateUser(parameters.ID, body.Name, body.Description)
	if err == database.UserNotFound {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		context.IndentedJSON(http.StatusOK, userData)
	}
}

// @Summary     Create user CRUD operation
// @Description Creates user in the Postgres database and returns created user as a JSON object
// @Tags        User
// @Accept      json
// @Produce     json
// @Param   	parameters body CreateUserData true "Create user parameters"
// @Failure     400 {object} APIError "Bad request"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /user [post]
func (h *Handler) createUser(context *gin.Context) {
	log.Println("Create user")

	var body CreateUserData
	if err := context.ShouldBindBodyWithJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := h.application.Database.CreateUser(body.Name, body.Description)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		context.IndentedJSON(http.StatusOK, userData)
	}
}

// @Summary     Delete user CRUD operation
// @Description Delete user from the Postgres database and returns deleted user as a JSON object
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Failure     400 {object} APIError "Bad request"
// @Failure     404 {object} APIError "User not found"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /user/{id} [delete]
func (h *Handler) deleteUser(context *gin.Context) {
	log.Println("Delete user")

	var parameters UserRequestParameters
	if err := context.ShouldBindUri(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := h.application.Database.DeleteUser(parameters.ID)
	if err == database.UserNotFound {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		context.IndentedJSON(http.StatusOK, userData)
	}
}
