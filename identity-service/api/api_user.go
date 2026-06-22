package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/alebik0/go-auth/identity-service/data"
	"github.com/gin-gonic/gin"
)

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
func (handler *Handler) ReadMyUser(context *gin.Context) {
	log.Printf("Read authentificated user")

	log.Printf("Check permissions")
	if !handler.isUser(context) {
		context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
		return
	}

	log.Printf("Load authorized user data")
	userID, err := handler.getUserID(context)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
		return
	}

	log.Printf("Read my user")
	query := "SELECT id, name, description FROM users WHERE id = $1;"
	userData := data.UserData{}
	err = handler.
		database.
		QueryRow(query, userID).
		Scan(&userData.ID, &userData.Name, &userData.Description)

	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		log.Printf("[ERROR] %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
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
func (handler *Handler) ReadUser(context *gin.Context) {
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
	userData := data.UserData{}
	err = handler.
		database.
		QueryRow(query, parameters.ID).
		Scan(&userData.ID, &userData.Name, &userData.Description)

	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		log.Printf("[ERROR] %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
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
// @Param   	parameters body UpdateUserRequest true "Update user request body parameters"
// @Param 		X-User-ID header string true "Authorized user ID"
// @Param 		X-User-Role header string true "Authorized user roles"
// @Failure     400 {object} APIError "Bad request"
// @Failure     403 {object} APIError "Forbidden"
// @Failure     404 {object} APIError "User not found"
// @Failure     500 {object} APIError "Internal server error"
// @Router      /api/v1/users/{id} [put]
func (handler *Handler) UpdateUser(context *gin.Context) {
	log.Printf("Update user")

	log.Printf("Load URI parameters")
	var parameters UserRequestParameters
	if err := context.ShouldBindUri(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Check permissions")
	if !handler.isAdmin(context) {
		if !handler.isUser(context) {
			context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
			return
		} else {
			userID, err := handler.getUserID(context)
			if err != nil {
				context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
				return
			}
			if userID != parameters.ID {
				context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
				return
			}
		}
	}

	log.Printf("Load body parameters")
	var body UpdateUserRequest
	if err := context.ShouldBindBodyWithJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Update user with id=%d", parameters.ID)
	query := "UPDATE users SET name = $1, description = $2 WHERE id=$3 RETURNING id, name, description;"
	userData := data.UserData{}
	err := handler.
		database.
		QueryRow(query, body.Name, body.Description, parameters.ID).
		Scan(&userData.ID, &userData.Name, &userData.Description)
	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		log.Printf("[ERROR] %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
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
func (handler *Handler) DeleteUser(context *gin.Context) {
	log.Printf("Delete user")

	log.Printf("Load URI parameters")
	var parameters UserRequestParameters
	if err := context.ShouldBindUri(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Check permissions")
	if !handler.isAdmin(context) {
		if !handler.isUser(context) {
			context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
			return
		} else {
			userID, err := handler.getUserID(context)
			if err != nil {
				context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
				return
			}
			if userID != parameters.ID {
				context.JSON(http.StatusForbidden, gin.H{"error": "Not enough permissions"})
				return
			}
		}
	}

	log.Printf("Delete user")
	query := "DELETE FROM users WHERE id=$1 RETURNING id, name, description;"
	userData := data.UserData{}
	err := handler.
		database.
		QueryRow(query, parameters.ID).
		Scan(&userData.ID, &userData.Name, &userData.Description)
	if errors.Is(err, sql.ErrNoRows) {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		log.Printf("[ERROR] %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	context.IndentedJSON(http.StatusOK, userData)
}
