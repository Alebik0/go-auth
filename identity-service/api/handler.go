package api

import (
	"database/sql"
	"errors"

	"github.com/redis/go-redis/v9"
)

type Handler struct {
	PostgresDatabase *sql.DB
	RedisDatabase    *redis.Client
	HmacSecret       []byte
}

// Dependency injection
func NewHandler(database *sql.DB, cache *redis.Client, hmacSecret []byte) (Handler, error) {
	return Handler{
		PostgresDatabase: database,
		RedisDatabase:    cache,
		HmacSecret:       hmacSecret,
	}, nil
}

func (handler *Handler) Close() error {
	return errors.Join(
		handler.PostgresDatabase.Close(),
		handler.RedisDatabase.Close(),
	)
}
