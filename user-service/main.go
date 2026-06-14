package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var app App

func main() {
	log.SetPrefix("[USER_SERVICE] ")

	var err error
	app, err = NewApp()
	if err != nil {
		log.Fatalf("Failed create app: %v", err)
	}
	defer app.Close()

	router := gin.Default()
	router.POST("/user", createUser)
	router.GET("/user/:id", readUser)
	router.PUT("/user/:id", updateUser)
	router.DELETE("/user/:id", deleteUser)

	log.Println("Running server on port 8080")
	router.Run("localhost:8080")
}
