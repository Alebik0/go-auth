package userservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	ApiTimeout = 1 * time.Second
)

type RemoteUserServiceAPI struct {
	ServiceHost string
	ServicePort string
}

func NewRemoteUserServiceAPI(host, port string) *RemoteUserServiceAPI {
	return &RemoteUserServiceAPI{
		ServiceHost: host,
		ServicePort: port,
	}
}

func (api *RemoteUserServiceAPI) CreateUser(name, description string) (UserData, error) {
	log.Printf("Prepare POST data")
	jsonData, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})
	if err != nil {
		return UserData{}, fmt.Errorf("failed to prepare post data: %w", err)
	}

	log.Printf("Prepare context")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		ApiTimeout,
	)
	defer cancel()

	log.Printf("Send POST /api/v1/user")
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s:%s/api/v1/user", api.ServiceHost, api.ServicePort),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to create post request: %w", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to send post request: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("[WARN] Failed to close request body")
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return UserData{}, ErrUserNotFound
	} else if resp.StatusCode != http.StatusOK {
		return UserData{}, fmt.Errorf("failed to send post request: status = %s (%d)", resp.Status, resp.StatusCode)
	}

	var response UserData
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to read body: %w", err)
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return UserData{}, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return response, nil
}

func (api *RemoteUserServiceAPI) ReadUser(id uint32) (UserData, error) {
	log.Printf("Prepare context")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		ApiTimeout,
	)
	defer cancel()

	log.Printf("Send GET /api/v1/user/{id}")
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s:%s/api/v1/%d", api.ServiceHost, api.ServicePort, id),
		http.NoBody,
	)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to create GET request: %w", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to send GET request: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("[WARN] Failed to close request body")
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return UserData{}, ErrUserNotFound
	} else if resp.StatusCode != http.StatusOK {
		return UserData{}, fmt.Errorf("failed to send GET request: status = %s (%d)", resp.Status, resp.StatusCode)
	}

	var response UserData
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to read body: %w", err)
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return UserData{}, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return response, nil
}

func (api *RemoteUserServiceAPI) UpdateUser(id uint32, name, description string) (UserData, error) {
	log.Printf("Prepare PUT data")
	jsonData, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})
	if err != nil {
		return UserData{}, fmt.Errorf("failed to prepare PUT data: %w", err)
	}

	log.Printf("Prepare context")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		ApiTimeout,
	)
	defer cancel()

	log.Printf("Send PUT /api/v1/user/{id}")
	req, err := http.NewRequestWithContext(
		ctx,
		"PUT",
		fmt.Sprintf("%s:%s/api/v1/%d", api.ServiceHost, api.ServicePort, id),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to create PUT request: %w", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to send PUT request: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("[WARN] Failed to close request body")
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return UserData{}, ErrUserNotFound
	} else if resp.StatusCode != http.StatusOK {
		return UserData{}, fmt.Errorf("failed to send post request: status = %s (%d)", resp.Status, resp.StatusCode)
	}

	var response UserData
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to read body: %w", err)
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return UserData{}, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return response, nil
}

func (api *RemoteUserServiceAPI) DeleteUser(id uint32) (UserData, error) {
	log.Printf("Prepare context")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		ApiTimeout,
	)
	defer cancel()

	log.Printf("Send DELETE /api/v1/user/{id}")
	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",
		fmt.Sprintf("%s:%s/api/v1/%d", api.ServiceHost, api.ServicePort, id),
		http.NoBody,
	)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to send DELETE request: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("[WARN] Failed to close request body")
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return UserData{}, ErrUserNotFound
	} else if resp.StatusCode != http.StatusOK {
		return UserData{}, fmt.Errorf("failed to send post request: status = %s (%d)", resp.Status, resp.StatusCode)
	}

	var response UserData
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to read body: %w", err)
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return UserData{}, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return response, nil
}

func (api *RemoteUserServiceAPI) Close() error {
	return nil
}
