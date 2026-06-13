package database

// CRUD user API
type DatabaseAPI interface {
	CreateUser(name, description string) (UserData, error)
	ReadUser(id uint32) (UserData, error)
	UpdateUser(id uint32, name, description string) (UserData, error)
	DeleteUser(id uint32) (UserData, error)
	Close() error
}

type UserData struct {
	ID          uint32
	Name        string
	Description string
}
