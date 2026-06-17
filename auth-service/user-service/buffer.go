package userservice

var counter uint32 = 1

type BufferDatabaseAPI struct {
	data map[uint32]UserData
}

func NewBufferDatabaseAPI() *BufferDatabaseAPI {
	return &BufferDatabaseAPI{
		data: make(map[uint32]UserData),
	}
}

func (api *BufferDatabaseAPI) CreateUser(name, description string) (UserData, error) {
	id := counter
	counter += 1
	userData := UserData{
		ID:          id,
		Name:        name,
		Description: description,
	}
	api.data[id] = userData

	return userData, nil
}

func (api *BufferDatabaseAPI) ReadUser(id uint32) (UserData, error) {
	userData, ok := api.data[id]
	if !ok {
		return userData, ErrUserNotFound
	}
	return userData, nil
}

func (api *BufferDatabaseAPI) UpdateUser(id uint32, name, description string) (UserData, error) {
	userData, ok := api.data[id]
	if !ok {
		return userData, ErrUserNotFound
	}
	newUserData := UserData{
		ID:          userData.ID,
		Name:        name,
		Description: description,
	}
	api.data[userData.ID] = newUserData

	return newUserData, nil
}

func (api *BufferDatabaseAPI) DeleteUser(id uint32) (UserData, error) {
	userData, ok := api.data[id]
	if !ok {
		return userData, ErrUserNotFound
	}

	delete(api.data, id)

	return userData, nil
}

func (api *BufferDatabaseAPI) Close() error {
	return nil
}
