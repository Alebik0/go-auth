package database

import (
	"errors"
)

// CRUD user API
type DatabaseAPI interface {
	CreateAuth(login, password_hash string, profile uint32) (AuthData, error)
	ReadAuthByID(id uint32) (AuthData, error)
	ReadAuthByLogin(login string) (AuthData, error)
	UpdateAuth(id uint32, passwordHash string) (AuthData, error)
	DeleteAuth(id uint32) (AuthData, error)
	Close() error
}

type AuthData struct {
	ID           uint32 `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
	ProfileID    uint32 `json:"profile_id"`
}

var (
	ErrAuthNotFound = errors.New("auth not found")
)
