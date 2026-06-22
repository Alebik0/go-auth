package api

import (
	"database/sql"
	"errors"

	"github.com/redis/go-redis/v9"
)

type Handler struct {
	database   *sql.DB
	cache      *redis.Client
	hmacSecret []byte
}

// Dependency injection
func NewHandler(database *sql.DB, cache *redis.Client, hmacSecret []byte) (Handler, error) {
	return Handler{
		database:   database,
		cache:      cache,
		hmacSecret: hmacSecret,
	}, nil
}

func (handler *Handler) Close() error {
	return errors.Join(
		handler.database.Close(),
		handler.cache.Close(),
	)
}
