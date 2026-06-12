package main

import (
	"log"
	"os"
	"time"

	"github.com/alebik0/go-auth/user-service/database"
)

func main() {
	api, err := database.NewPostgresAPI("localhost", 5432, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to load database API: %v", err)
	}

	user, err := api.ReadUser(0)
	log.Printf("Loading user: %v, %v", user, err)
}
