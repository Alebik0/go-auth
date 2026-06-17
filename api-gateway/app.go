package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/alebik0/go-auth/auth-service/database"
	userservice "github.com/alebik0/go-auth/auth-service/user-service"
)

type Handler struct {
	AuthAPI database.DatabaseAPI
	UserAPI userservice.DatabaseAPI
}

// Dependency injection
func NewHaldler() (Handler, error) {
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

	// redisHost := os.Getenv("REDIS_HOST")
	// if redisHost == "" {
	// 	return Handler{}, fmt.Errorf("REDIS_HOST is mandatory environment variable")
	// }
	// redisPort := os.Getenv("REDIS_PORT")
	// if redisPort == "" {
	// 	return Handler{}, fmt.Errorf("REDIS_PORT is mandatory environment variable")
	// }
	// redisPassword := os.Getenv("REDIS_PASSWORD")
	// if redisPassword == "" {
	// 	return Handler{}, fmt.Errorf("REDIS_PASSWORD is mandatory environment variable")
	// }
	// redisDatabase := os.Getenv("REDIS_DATABASE")
	// if redisDatabase == "" {
	// 	return Handler{}, fmt.Errorf("REDIS_DATABASE is mandatory environment variable")
	// }
	// redisDatabaseInt, err := strconv.ParseInt(redisDatabase, 10, 32)
	// if err != nil {
	// 	return Handler{}, fmt.Errorf("REDIS_DATABASE must be an integer")
	// }

	hmacSecret := os.Getenv("AUTH_SERVICE_HMAC_SECRET")
	if hmacSecret == "" {
		return Handler{}, fmt.Errorf("AUTH_SERVICE_HMAC_SECRET is mandatory environment variable")
	}

	return Handler{
		AuthAPI: authApi,
		UserAPI: userApi,
	}, nil
}

func (handler *Handler) Close() error {
	return errors.Join(
		handler.AuthAPI.Close(),
		handler.UserAPI.Close(),
	)
}
