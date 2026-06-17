package jwt

import (
	"context"
	"errors"
	"time"
)

type JwtDatabaseAPI interface {
	Save(ctx context.Context, token, userID string, expiration time.Time) error
	Get(ctx context.Context, token string) (string, error)
	Revert(ctx context.Context, token string) error
	Close() error
}

var (
	ErrJwtNotFound = errors.New("jwt not found")
)
