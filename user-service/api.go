package main

import (
	"log"
	"net/http"

	"github.com/alebik0/go-auth/user-service/database"
	"github.com/gin-gonic/gin"
)

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

func readUser(context *gin.Context) {
	log.Println("Read user")

	var parameters UserRequestParameters
	if err := context.ShouldBindUri(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := app.Database.ReadUser(parameters.ID)
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

func updateUser(context *gin.Context) {
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

	userData, err := app.Database.UpdateUser(parameters.ID, body.Name, body.Description)
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

func createUser(context *gin.Context) {
	log.Println("Create user")

	var body UpdateUserData
	if err := context.ShouldBindBodyWithJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := app.Database.CreateUser(body.Name, body.Description)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		context.IndentedJSON(http.StatusOK, userData)
	}
}

func deleteUser(context *gin.Context) {
	log.Println("Delete user")

	var parameters UserRequestParameters
	if err := context.ShouldBindUri(&parameters); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := app.Database.DeleteUser(parameters.ID)
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
