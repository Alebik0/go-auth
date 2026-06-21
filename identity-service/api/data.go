package api

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type UserRequestParameters struct {
	ID uint32 `uri:"id"`
}

type UpdateUserRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type APIError struct {
	Error string `json:"error" example:"just a random internal error"`
}
