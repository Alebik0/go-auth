package jwt

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisJwtDatabaseAPI struct {
	client *redis.Client
}

func NewRedisJwtDatabaseAPI(addr, password string, database int) *RedisJwtDatabaseAPI {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	})

	return &RedisJwtDatabaseAPI{client: rdb}
}

func (api *RedisJwtDatabaseAPI) Save(ctx context.Context, token, userID string, expiration time.Time) error {
	log.Printf("Save token: key=%s, value=%s", token, userID)
	err := api.client.Set(ctx, "token:"+token, userID, time.Until(expiration)).Err()
	if err != nil {
		return fmt.Errorf("failed to set to redis: %w", err)
	}
	return nil
}

func (api *RedisJwtDatabaseAPI) Get(ctx context.Context, token string) (string, error) {
	log.Printf("Get token: key=%s", token)
	result, err := api.client.Get(ctx, "token:"+token).Result()
	if err == redis.Nil {
		return "", ErrJwtNotFound
	} else if err != nil {
		return "", fmt.Errorf("failed to get from redis: %w", err)
	}
	return result, nil
}

func (api *RedisJwtDatabaseAPI) Revert(ctx context.Context, token string) error {
	log.Printf("Revert token: key=%s", token)
	err := api.client.Del(ctx, "token:"+token).Err()
	if err != nil {
		return fmt.Errorf("failed to revoke token redis: %w", err)
	}
	return nil
}

func (api *RedisJwtDatabaseAPI) Close() error {
	return api.client.Close()
}
