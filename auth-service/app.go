package main

import (
	"errors"

	"github.com/alebik0/go-auth/auth-service/database"
	userservice "github.com/alebik0/go-auth/auth-service/user-service"
)

type App struct {
	AuthAPI database.DatabaseAPI
	UserAPI userservice.DatabaseAPI
}

func (app *App) Close() error {
	err := errors.Join(
		app.AuthAPI.Close(),
		app.UserAPI.Close(),
	)

	return err
}

// Dependency injection logic
func NewApp() (App, error) {
	// host := os.Getenv("POSTGRES_HOST")
	// if host == "" {
	// 	return App{}, fmt.Errorf("POSTGRES_HOST is mandatory environment variable")
	// }
	// port := os.Getenv("POSTGRES_PORT")
	// if host == "" {
	// 	return App{}, fmt.Errorf("POSTGRES_PORT is mandatory environment variable")
	// }
	// portInt, err := strconv.ParseUint(port, 10, 16)
	// if err != nil {
	// 	return App{}, fmt.Errorf("POSTGRES_PORT must be a uint16 number")
	// }
	// user := os.Getenv("POSTGRES_USER")
	// if user == "" {
	// 	return App{}, fmt.Errorf("POSTGRES_USER is mandatory environment variable")
	// }
	// password := os.Getenv("POSTGRES_PASSWORD")
	// if password == "" {
	// 	return App{}, fmt.Errorf("POSTGRES_PASSWORD is mandatory environment variable")
	// }
	// dbName := os.Getenv("POSTGRES_DB")
	// if dbName == "" {
	// 	return App{}, fmt.Errorf("POSTGRES_DB is mandatory environment variable")
	// }

	// authApi, err := database.NewPostgresAPI(host, uint16(portInt), user, password, dbName, 5*time.Second)
	// if err != nil {
	// 	return App{}, fmt.Errorf("failed to load database API: %v", err)
	// }
	authApi := database.NewBufferDatabaseAPI()

	// userServiceHost := os.Getenv("USER_SERVICE_HOST")
	// if host == "" {
	// 	return App{}, fmt.Errorf("USER_SERVICE_HOST is mandatory environment variable")
	// }
	// userServicePort := os.Getenv("USER_SERVICE_PORT")
	// if host == "" {
	// 	return App{}, fmt.Errorf("USER_SERVICE_PORT is mandatory environment variable")
	// }

	// userApi := userservice.NewRemoteUserServiceAPI(userServiceHost, userServicePort)
	userApi := userservice.NewBufferDatabaseAPI()

	return App{
		AuthAPI: authApi,
		UserAPI: userApi,
	}, nil
}
