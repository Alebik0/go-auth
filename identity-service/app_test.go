package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/alebik0/go-auth/identity-service/api"
	"github.com/alebik0/go-auth/identity-service/data"
	"github.com/alebik0/go-auth/identity-service/jwt"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4/testutils/assert"

	"github.com/go-openapi/testify/v2/require"
	_ "github.com/lib/pq" // To register the driver.

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Logf("Prepare data")
	hmacSecret := "qweqeqwioueqiewuio"
	login := "username"
	password := "passwd"

	t.Logf("Create postgres API mock")
	database, mock, err := sqlmock.New()
	require.NoError(t, err)

	t.Logf("Create in-memory JWT API")
	jwtApi := jwt.NewBufferJwtDatabaseAPI()

	t.Logf("Create handler")
	handler, err := api.NewHandler(
		database,
		jwtApi,
		[]byte(hmacSecret),
	)
	assert.NoError(t, err)

	t.Logf("Setup router")
	router := SetupRouter(handler)
	globalCookies := make([]*http.Cookie, 0)
	authorization := ""

	{
		t.Logf("Prepare mock")
		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					"SELECT id, login, password_hash, profile_id FROM auth WHERE login = $1;",
				),
			).
			WithArgs(login).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectBegin()
		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					"INSERT INTO users (name, description) VALUES ($1, $2) RETURNING id;",
				),
			).
			WithArgs(login, "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					"INSERT INTO auth (login, password_hash, profile_id) VALUES ($1, $2, $3) RETURNING id;",
				),
			).
			WithArgs(login, sqlmock.AnyArg(), 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		t.Logf("POST /api/v1/auth/register")
		body := strings.NewReader(fmt.Sprintf(`{"login":"%s", "password":"%s"}`, login, password))
		request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
		responseWriter := httptest.NewRecorder()

		router.ServeHTTP(responseWriter, request)

		if responseWriter.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", responseWriter.Code)
		}

		var response api.RegisterResponse
		err = json.Unmarshal(responseWriter.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		cookies := responseWriter.Result().Cookies()

		if len(cookies) != 1 {
			t.Fatal("expected one cookie")
		}

		if cookies[0].Name != "refresh_token" {
			t.Fatal("wrong cookie")
		}

		globalCookies = append(globalCookies, cookies...)
		authorization = response.AccessToken
	}

	{
		t.Logf("Prepare mock")
		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					"SELECT id, name, description FROM users WHERE id = $1;",
				),
			).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, login, ""))

		t.Logf("GET /api/v1/users/my")
		body := http.NoBody
		request := httptest.NewRequest(http.MethodGet, "/api/v1/users/my", body)
		for _, c := range globalCookies {
			request.AddCookie(c)
		}
		request.Header.Add("Authorization", "Bearer "+authorization)

		responseWriter := httptest.NewRecorder()

		router.ServeHTTP(responseWriter, request)

		if responseWriter.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body: %s", responseWriter.Code, responseWriter.Body.Bytes())
		}

		var response data.UserData
		err = json.Unmarshal(responseWriter.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		expected := data.UserData{
			ID:          1,
			Name:        login,
			Description: "",
		}

		if !reflect.DeepEqual(response, expected) {
			t.Fatalf("got %+v want %+v", response, expected)
		}
	}
}
