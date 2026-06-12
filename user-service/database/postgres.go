package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq" // To register the driver.
)

type PostgresAPI struct {
	database *sql.DB
}

func prepareDatabase(db *sql.DB) error {
	log.Println("Prepare database: create tables if not exists")

	query := `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(64) NOT NULL UNIQUE,
    password_hash VARCHAR(128) NOT NULL,
    name VARCHAR(64) NOT NULL,
	description VARCHAR(1024) NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
	_, err := db.Query(query)

	if err != nil {
		return fmt.Errorf("failed run initiation query: %v", err)
	}

	return nil
}

func NewPostgresAPI(host string, port uint16, user, password, database string, connectTimeout time.Duration) (PostgresAPI, error) {
	log.Println("Create postgres API")

	cfg := pq.Config{
		Host:           host,
		Port:           port,
		User:           user,
		Password:       password,
		Database:       database,
		ConnectTimeout: connectTimeout,
		SSLMode:        pq.SSLModeDisable,
	}

	c, err := pq.NewConnectorConfig(cfg)
	if err != nil {
		return PostgresAPI{}, fmt.Errorf("failed create connection config: %v", err)
	}

	db := sql.OpenDB(c)

	log.Println("Ping postgres database")
	err = db.Ping()
	if err != nil {
		db.Close()
		return PostgresAPI{}, fmt.Errorf("failed ping server: %v", err)
	}

	err = prepareDatabase(db)
	if err != nil {
		db.Close()
		return PostgresAPI{}, fmt.Errorf("failed prepare database: %v", err)
	}

	return PostgresAPI{database: db}, nil
}

func (api *PostgresAPI) Close() error {
	return api.database.Close()
}

func (api *PostgresAPI) CreateUser(user UserData) error {
	return nil
}

func (api *PostgresAPI) ReadUser(id uint32) (UserData, error) {
	var userData UserData
	err := api.database.QueryRow("SELECT * FROM users WHERE id=$1", id).Scan(&userData.ID, &userData.Login, &userData.PasswordHash, &userData.Name, &userData.Description)
	if err == sql.ErrNoRows {
		return userData, fmt.Errorf("user not found")
	} else if err != nil {
		return userData, fmt.Errorf("failed to scan row: %v", err)
	} else {
		return userData, nil
	}
}

func (api *PostgresAPI) UpdateUser(id uint32, name, description *string) error {
	return nil
}

func (api *PostgresAPI) DeleteUser(id uint32) (UserData, error) {
	return UserData{}, nil
}
