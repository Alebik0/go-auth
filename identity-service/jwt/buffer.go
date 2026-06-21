package jwt

import (
	"context"
	"log"
	"time"
)

type BufferJwtRow struct {
	Value string
	Time  time.Time
}

type BufferJwtDatabaseAPI struct {
	data map[string]BufferJwtRow
}

func NewBufferJwtDatabaseAPI() *BufferJwtDatabaseAPI {
	return &BufferJwtDatabaseAPI{data: make(map[string]BufferJwtRow)}
}

func (api *BufferJwtDatabaseAPI) Save(ctx context.Context, token, userID string, expiration time.Time) error {
	log.Printf("Save token: key=%s, value=%s", token, userID)
	api.data["token:"+token] = BufferJwtRow{
		Value: userID,
		Time:  expiration,
	}
	return nil
}

func (api *BufferJwtDatabaseAPI) Get(ctx context.Context, token string) (string, error) {
	log.Printf("Get token: key=%s", token)
	result, ok := api.data["token:"+token]
	if !ok {
		return "", ErrJwtNotFound
	} else if result.Time.Before(time.Now()) {
		delete(api.data, "token:"+token)
		return "", ErrJwtNotFound
	}
	return result.Value, nil
}

func (api *BufferJwtDatabaseAPI) Revert(ctx context.Context, token string) error {
	log.Printf("Revert token: key=%s", token)
	delete(api.data, "token:"+token)
	return nil
}

func (api *BufferJwtDatabaseAPI) Close() error {
	return nil
}
