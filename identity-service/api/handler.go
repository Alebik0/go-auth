package api

import (
	"database/sql"
	"errors"

	"github.com/alebik0/go-auth/identity-service/jwt"
)

type Handler struct {
	jwtAPI     jwt.JwtDatabaseAPI
	database   *sql.DB
	hmacSecret []byte
}

// Dependency injection
func NewHandler(database *sql.DB, jwtApi jwt.JwtDatabaseAPI, hmacSecret []byte) (Handler, error) {
	return Handler{
		database:   database,
		jwtAPI:     jwtApi,
		hmacSecret: hmacSecret,
	}, nil
}

func (handler *Handler) Close() error {
	return errors.Join(
		handler.database.Close(),
		handler.jwtAPI.Close(),
	)
}
