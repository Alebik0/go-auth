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
		err := db.Close()
		if err != nil {
			log.Printf("[WARN] Failed to close database: %v", err)
		}

		return nil, fmt.Errorf("failed ping server: %v", err)
	}

	log.Println("Prepare postgres database")
	err = prepareDatabase(db)
	if err != nil {
		err := db.Close()
		if err != nil {
			log.Printf("[WARN] Failed to close database: %v", err)
		}

		return nil, fmt.Errorf("failed prepare database: %v", err)
	}

	log.Println("Created postgres API")

	return &PostgresAPI{database: db}, nil
}

func (api *PostgresAPI) CreateAuth(login, passwordHash string, profileID uint32) (AuthData, error) {
	query := "INSERT INTO auth (login, password_hash, profile_id) VALUES ($1, $2, $3) RETURNING id;"
	authData := AuthData{
		Login:        login,
		PasswordHash: passwordHash,
		ProfileID:    profileID,
	}
	err := api.database.QueryRow(query, login, passwordHash, profileID).Scan(&authData.ID)
	if err != nil {
		return authData, fmt.Errorf("failed to scan row: %v", err)
	}
	return authData, nil
}

func (api *PostgresAPI) ReadAuthByID(id uint32) (AuthData, error) {
	query := "SELECT id, login, password_hash, profile_id FROM auth WHERE id = $1;"
	authData := AuthData{}
	err := api.database.QueryRow(query, id).Scan(&authData.ID, &authData.Login, &authData.PasswordHash, &authData.ProfileID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authData, ErrAuthNotFound
		}
		return authData, fmt.Errorf("failed to scan row: %v", err)
	}
	return authData, nil
}

func (api *PostgresAPI) ReadAuthByLogin(login string) (AuthData, error) {
	query := "SELECT id, login, password_hash, profile_id FROM auth WHERE login = $1;"
	authData := AuthData{}
	err := api.database.QueryRow(query, login).Scan(&authData.ID, &authData.Login, &authData.PasswordHash, &authData.ProfileID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authData, ErrAuthNotFound
		}
		return authData, fmt.Errorf("failed to scan row: %v", err)
	}
	return authData, nil
}

func (api *PostgresAPI) UpdateAuth(id uint32, passwordHash string) (AuthData, error) {
	query := "UPDATE auth SET password_hash = $1 WHERE id=$2 RETURNING id, login, password_hash, profile_id;"
	authData := AuthData{}
	err := api.database.QueryRow(query, passwordHash, id).Scan(&authData.ID, &authData.Login, &authData.PasswordHash, &authData.ProfileID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authData, ErrAuthNotFound
		}
		return authData, fmt.Errorf("failed to scan row: %v", err)
	}
	return authData, nil
}

func (api *PostgresAPI) DeleteAuth(id uint32) (AuthData, error) {
	query := "DELETE FROM auth WHERE id=$1 RETURNING id, login, password_hash, profile_id;"
	authData := AuthData{}
	err := api.database.QueryRow(query, id).Scan(&authData.ID, &authData.Login, &authData.PasswordHash, &authData.ProfileID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authData, ErrAuthNotFound
		}
		return authData, fmt.Errorf("failed to scan row: %v", err)
	}
	return authData, nil
}

func (api *PostgresAPI) Close() error {
	return api.database.Close()
}
