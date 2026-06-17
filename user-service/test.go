package main

import (
	"log"

	"github.com/alebik0/go-auth/user-service/database"
)

func testAPI() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed create app: %v", err)
	}
	defer app.Close()

	log.Printf("Create user")
	user, err := app.Database.CreateUser("user", "description")
	if err != nil {
		log.Fatalf("Expected nil, but got: %v", err)
	}
	if user.Name != "user" {
		log.Fatalf("Expected 'user', but got: %v", user.Name)
	}
	if user.Description != "description" {
		log.Fatalf("Expected 'description', but got: %v", user.Description)
	}

	log.Printf("Created user #%d", user.ID)

	log.Printf("Read user")
	readUser, err := app.Database.ReadUser(user.ID)
	if err != nil {
		log.Fatalf("Expected nil, but got: %v", err)
	}
	if readUser.Name != "user" {
		log.Fatalf("Expected 'user', but got: %v", readUser.Name)
	}
	if readUser.Description != "description" {
		log.Fatalf("Expected 'description', but got: %v", readUser.Description)
	}

	log.Printf("Update user")
	updatedUser, err := app.Database.UpdateUser(user.ID, "user", "description#2")
	if err != nil {
		log.Fatalf("Expected nil, but got: %v", err)
	}
	if updatedUser.Name != "user" {
		log.Fatalf("Expected 'user', but got: %v", updatedUser.Name)
	}
	if updatedUser.Description != "description#2" {
		log.Fatalf("Expected 'description#2', but got: %v", updatedUser.Description)
	}

	log.Printf("Read user (again)")
	readUser2, err := app.Database.ReadUser(user.ID)
	if err != nil {
		log.Fatalf("Expected nil, but got: %v", err)
	}
	if readUser2.Name != "user" {
		log.Fatalf("Expected 'user', but got: %v", readUser2.Name)
	}
	if readUser2.Description != "description#2" {
		log.Fatalf("Expected 'description#2', but got: %v", readUser2.Description)
	}

	log.Printf("Delete user")
	deletedUser, err := app.Database.DeleteUser(user.ID)
	if err != nil {
		log.Fatalf("Expected nil, but got: %v", err)
	}
	if deletedUser.Name != "user" {
		log.Fatalf("Expected 'user', but got: %v", deletedUser.Name)
	}
	if deletedUser.Description != "description#2" {
		log.Fatalf("Expected 'description#2', but got: %v", deletedUser.Description)
	}

	log.Printf("Read user (again x2)")
	_, err = app.Database.ReadUser(user.ID)
	if err != database.ErrUserNotFound {
		log.Fatalf("Expected UserNotFound, but got: %v", err)
	}

	log.Printf("Success")
}
