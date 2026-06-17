package database

var counter uint32 = 1

type BufferDatabaseAPI struct {
	data map[uint32]AuthData
}

func NewBufferDatabaseAPI() *BufferDatabaseAPI {
	return &BufferDatabaseAPI{
		data: make(map[uint32]AuthData),
	}
}

func (api *BufferDatabaseAPI) CreateAuth(login, passwordHash string, profileID uint32) (AuthData, error) {
	id := counter
	counter += 1
	authData := AuthData{
		ID:           id,
		Login:        login,
		PasswordHash: passwordHash,
		ProfileID:    profileID,
	}
	api.data[id] = authData

	return authData, nil
}

func (api *BufferDatabaseAPI) ReadAuthByID(id uint32) (AuthData, error) {
	authData, ok := api.data[id]
	if !ok {
		return authData, ErrAuthNotFound
	}
	return authData, nil
}

func (api *BufferDatabaseAPI) ReadAuthByLogin(login string) (AuthData, error) {
	for _, authData := range api.data {
		if authData.Login == login {
			return authData, nil
		}
	}
	return AuthData{}, ErrAuthNotFound
}

func (api *BufferDatabaseAPI) UpdateAuth(id uint32, passwordHash string) (AuthData, error) {
	authData, ok := api.data[id]
	if !ok {
		return authData, ErrAuthNotFound
	}
	newAuthData := AuthData{
		ID:           authData.ID,
		Login:        authData.Login,
		PasswordHash: passwordHash,
		ProfileID:    authData.ProfileID,
	}
	api.data[authData.ID] = newAuthData

	return newAuthData, nil
}

func (api *BufferDatabaseAPI) DeleteAuth(id uint32) (AuthData, error) {
	authData, ok := api.data[id]
	if !ok {
		return authData, ErrAuthNotFound
	}

	delete(api.data, id)

	return authData, nil
}

func (api *BufferDatabaseAPI) Close() error {
	return nil
}
