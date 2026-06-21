package data

type AuthData struct {
	ID           uint32 `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
	ProfileID    uint32 `json:"profile_id"`
}

type UserData struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
