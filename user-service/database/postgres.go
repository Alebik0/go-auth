package database

import (
	"database/sql"
	"errors"
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
DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
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

func NewPostgresAPI(host string, port uint16, user, password, database string, connectTimeout time.Duration) (*PostgresAPI, error) {
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
		return nil, fmt.Errorf("failed create connection config: %v", err)
	}

	db := sql.OpenDB(c)

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Ping postgres database")
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed ping server: %v", err)
	}

	log.Println("Prepare postgres database")
	err = prepareDatabase(db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed prepare database: %v", err)
	}

	log.Println("Created postgres API")

	return &PostgresAPI{database: db}, nil
}

func (api *PostgresAPI) Close() error {
	return api.database.Close()
}

func (api *PostgresAPI) CreateUser(name, description string) (UserData, error) {
	query := "INSERT INTO users (name, description) VALUES ($1, $2) RETURNING id;"
	userData := UserData{
		Name:        name,
		Description: description,
	}
	err := api.database.QueryRow(query, name, description).Scan(&userData.ID)
	if err != nil {
		return userData, fmt.Errorf("failed to scan row: %v", err)
	}
	return userData, nil
}

func (api *PostgresAPI) ReadUser(id uint32) (UserData, error) {
	query := "SELECT id, name, description FROM users WHERE id = $1;"
	userData := UserData{}
	err := api.database.QueryRow(query, id).Scan(&userData.ID, &userData.Name, &userData.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userData, UserNotFound
		}
		return userData, fmt.Errorf("failed to scan row: %v", err)
	}
	return userData, nil
}

func (api *PostgresAPI) UpdateUser(id uint32, name, description string) (UserData, error) {
	query := "UPDATE users SET name = $1, description = $2 WHERE id=$3 RETURNING id, name, description;"
	userData := UserData{}
	err := api.database.QueryRow(query, name, description, id).Scan(&userData.ID, &userData.Name, &userData.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userData, UserNotFound
		}
		return userData, fmt.Errorf("failed to scan row: %v", err)
	}
	return userData, nil
}

func (api *PostgresAPI) DeleteUser(id uint32) (UserData, error) {
	query := "DELETE FROM users WHERE id=$1 RETURNING id, name, description;"
	userData := UserData{}
	err := api.database.QueryRow(query, id).Scan(&userData.ID, &userData.Name, &userData.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userData, UserNotFound
		}
		return userData, fmt.Errorf("failed to scan row: %v", err)
	}
	return userData, nil
}
