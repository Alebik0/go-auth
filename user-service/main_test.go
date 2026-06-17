package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alebik0/go-auth/user-service/database"
	"github.com/go-openapi/testify/v2/assert"
)

func TestReadNotFount(t *testing.T) {
	t.Logf("Init dependencies")
	app := App{Database: database.NewBufferDatabaseAPI()}
	handler := NewHaldler(app)
	router := SetupRouter(handler)

	t.Logf("Send request")
	req, err := http.NewRequest("GET", "/api/v1/users/1", http.NoBody)
	assert.NoError(t, err)

	responseWriter := httptest.NewRecorder()
	router.ServeHTTP(responseWriter, req)

	t.Logf("Check response")
	assert.Equal(t, http.StatusNotFound, responseWriter.Code)
}

func TestCreateAndRead(t *testing.T) {
	t.Logf("Init dependencies")
	app := App{Database: database.NewBufferDatabaseAPI()}
	handler := NewHaldler(app)
	router := SetupRouter(handler)

	{
		t.Logf("Prepare POST data")
		jsonData, err := json.Marshal(map[string]string{
			"name":        "hello",
			"description": "world",
		})
		assert.NoError(t, err)

		t.Logf("Send POST /api/v1/users")
		req, err := http.NewRequest(
			"POST",
			"/api/v1/users",
			bytes.NewBuffer(jsonData),
		)
		assert.NoError(t, err)

		responseWriter := httptest.NewRecorder()
		router.ServeHTTP(responseWriter, req)

		t.Logf("Check response")
		t.Logf("Got status code: %d", responseWriter.Code)
		t.Logf("Got response body: %s", responseWriter.Body.String())
		assert.Equal(t, http.StatusOK, responseWriter.Code)

		expected := database.UserData{
			ID:          1,
			Name:        "hello",
			Description: "world",
		}
		var response database.UserData
		err = json.Unmarshal(responseWriter.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expected, response)
	}

	{
		t.Logf("Prepare PUT data")
		jsonData, err := json.Marshal(map[string]string{
			"name":        "hello",
			"description": "new world",
		})
		assert.NoError(t, err)

		t.Logf("Send POST /api/v1/users/{id}")
		req, err := http.NewRequest(
			"PUT",
			"/api/v1/users/1",
			bytes.NewBuffer(jsonData),
		)
		assert.NoError(t, err)

		responseWriter := httptest.NewRecorder()
		router.ServeHTTP(responseWriter, req)

		t.Logf("Check response")
		t.Logf("Got status code: %d", responseWriter.Code)
		t.Logf("Got response body: %s", responseWriter.Body.String())
		assert.Equal(t, http.StatusOK, responseWriter.Code)

		expected := database.UserData{
			ID:          1,
			Name:        "hello",
			Description: "new world",
		}
		var response database.UserData
		err = json.Unmarshal(responseWriter.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expected, response)
	}

	{
		t.Logf("Send GET /api/v1/users/{id}")
		req, err := http.NewRequest(
			"GET",
			"/api/v1/users/1",
			http.NoBody,
		)
		assert.NoError(t, err)

		responseWriter := httptest.NewRecorder()
		router.ServeHTTP(responseWriter, req)

		t.Logf("Check response")
		t.Logf("Got status code: %d", responseWriter.Code)
		t.Logf("Got response body: %s", responseWriter.Body.String())
		assert.Equal(t, http.StatusOK, responseWriter.Code)

		expected := database.UserData{
			ID:          1,
			Name:        "hello",
			Description: "new world",
		}
		var response database.UserData
		err = json.Unmarshal(responseWriter.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expected, response)
	}

	{
		t.Logf("Send DELETE /api/v1/users/{id}")
		req, err := http.NewRequest(
			"DELETE",
			"/api/v1/users/1",
			http.NoBody,
		)
		assert.NoError(t, err)

		responseWriter := httptest.NewRecorder()
		router.ServeHTTP(responseWriter, req)

		t.Logf("Check response")
		t.Logf("Got status code: %d", responseWriter.Code)
		t.Logf("Got response body: %s", responseWriter.Body.String())
		assert.Equal(t, http.StatusOK, responseWriter.Code)

		expected := database.UserData{
			ID:          1,
			Name:        "hello",
			Description: "new world",
		}
		var response database.UserData
		err = json.Unmarshal(responseWriter.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expected, response)
	}

	{
		t.Logf("Send GET /api/v1/users/{id}")
		req, err := http.NewRequest(
			"GET",
			"/api/v1/users/1",
			http.NoBody,
		)
		assert.NoError(t, err)

		responseWriter := httptest.NewRecorder()
		router.ServeHTTP(responseWriter, req)

		t.Logf("Check response")
		t.Logf("Got status code: %d", responseWriter.Code)
		t.Logf("Got response body: %s", responseWriter.Body.String())
		assert.Equal(t, http.StatusNotFound, responseWriter.Code)
	}
}
