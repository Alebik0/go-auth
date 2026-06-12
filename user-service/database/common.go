package database

type DatabaseAPI interface {
	CreateUser(user UserData) error
	ReadUser(id uint32) (UserData, error)
	UpdateUser(id uint32, name, description *string) error
	DeleteUser(id uint32) (UserData, error)
	Close() error
}

type UserData struct {
	ID           uint32
	Login        string
	PasswordHash string
	Name         string
	Description  string
}
