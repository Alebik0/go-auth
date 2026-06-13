package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/alebik0/go-auth/user-service/database"
)

type App struct {
	Database database.DatabaseAPI
}

func (app *App) Close() error {
	err := errors.Join(
		app.Database.Close(),
	)

	return err
}

// Dependency injection logic
func NewApp() (App, error) {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		return App{}, fmt.Errorf("POSTGRES_HOST is mandatory environment variable")
	}
	port := os.Getenv("POSTGRES_PORT")
	if host == "" {
		return App{}, fmt.Errorf("POSTGRES_PORT is mandatory environment variable")
	}
	portInt, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return App{}, fmt.Errorf("POSTGRES_PORT must be a uint16 number")
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		return App{}, fmt.Errorf("POSTGRES_USER is mandatory environment variable")
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		return App{}, fmt.Errorf("POSTGRES_PASSWORD is mandatory environment variable")
	}
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		return App{}, fmt.Errorf("POSTGRES_DB is mandatory environment variable")
	}

	api, err := database.NewPostgresAPI(host, uint16(portInt), user, password, dbName, 5*time.Second)
	if err != nil {
		return App{}, fmt.Errorf("failed to load database API: %v", err)
	}

	return App{Database: api}, nil
}
