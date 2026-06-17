package database

import "fmt"

// CRUD user API
type DatabaseAPI interface {
	CreateUser(name, description string) (UserData, error)
	ReadUser(id uint32) (UserData, error)
	UpdateUser(id uint32, name, description string) (UserData, error)
	DeleteUser(id uint32) (UserData, error)
	Close() error
}

type UserData struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var (
	UserNotFound = fmt.Errorf("user not found")
)
