package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/alebik0/go-auth/identity-service/jwt"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // To register the driver.
)

type Handler struct {
	JwtAPI     jwt.JwtDatabaseAPI
	Database   *sql.DB
	hmacSecret []byte
}

// Dependency injection
func NewHaldler() (Handler, error) {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		return Handler{}, fmt.Errorf("POSTGRES_HOST is mandatory environment variable")
	}
	port := os.Getenv("POSTGRES_PORT")
	if host == "" {
		return Handler{}, fmt.Errorf("POSTGRES_PORT is mandatory environment variable")
	}
	portInt, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return Handler{}, fmt.Errorf("POSTGRES_PORT must be a uint16 number")
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		return Handler{}, fmt.Errorf("POSTGRES_USER is mandatory environment variable")
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		return Handler{}, fmt.Errorf("POSTGRES_PASSWORD is mandatory environment variable")
	}
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		return Handler{}, fmt.Errorf("POSTGRES_DB is mandatory environment variable")
	}

	log.Println("Create postgres API")

	cfg := pq.Config{
		Host:           host,
		Port:           uint16(portInt),
		User:           user,
		Password:       password,
		Database:       dbName,
		ConnectTimeout: 5 * time.Second,
		SSLMode:        pq.SSLModeDisable,
	}

	c, err := pq.NewConnectorConfig(cfg)
	if err != nil {
		return Handler{}, fmt.Errorf("failed create connection config: %v", err)
	}

	db := sql.OpenDB(c)

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Ping postgres database")
	err = db.Ping()
	if err != nil {
		err := db.Close()
		if err != nil {
			log.Printf("[WARN] Failed to close database: %v", err)
		}

		return Handler{}, fmt.Errorf("failed ping server: %v", err)
	}

	log.Println("Prepare postgres database")

	query := `
DROP TABLE IF EXISTS auth;
CREATE TABLE IF NOT EXISTS auth (
    id SERIAL PRIMARY KEY,
	login VARCHAR(64) NOT NULL,
	password_hash VARCHAR(128) NOT NULL,
    profile_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_auth_profile
        FOREIGN KEY (profile_id)
        REFERENCES users(id)
);`
	_, err = db.Query(query)
	if err != nil {
		err := db.Close()
		if err != nil {
			log.Printf("[WARN] Failed to close database: %v", err)
		}

		return Handler{}, fmt.Errorf("failed prepare database: %v", err)
	}

	log.Println("Created postgres API")

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		return Handler{}, fmt.Errorf("REDIS_HOST is mandatory environment variable")
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		return Handler{}, fmt.Errorf("REDIS_PORT is mandatory environment variable")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		return Handler{}, fmt.Errorf("REDIS_PASSWORD is mandatory environment variable")
	}
	redisDatabase := os.Getenv("REDIS_DATABASE")
	if redisDatabase == "" {
		return Handler{}, fmt.Errorf("REDIS_DATABASE is mandatory environment variable")
	}
	redisDatabaseInt, err := strconv.ParseInt(redisDatabase, 10, 32)
	if err != nil {
		return Handler{}, fmt.Errorf("REDIS_DATABASE must be an integer")
	}

	jwtApi := jwt.NewRedisJwtDatabaseAPI(
		fmt.Sprintf("%s:%s", redisHost, redisPort),
		redisPassword,
		int(redisDatabaseInt),
	)
	// jwtApi := jwt.NewBufferJwtDatabaseAPI()

	hmacSecret := os.Getenv("IDENTITY_SERVICE_HMAC_SECRET")
	if hmacSecret == "" {
		return Handler{}, fmt.Errorf("IDENTITY_SERVICE_HMAC_SECRET is mandatory environment variable")
	}

	return Handler{
		Database:   db,
		JwtAPI:     jwtApi,
		hmacSecret: []byte(hmacSecret),
	}, nil
}

func (handler *Handler) Close() error {
	return errors.Join(
		handler.Database.Close(),
		handler.JwtAPI.Close(),
	)
}
