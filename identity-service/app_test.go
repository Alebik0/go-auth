package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alebik0/go-auth/identity-service/api"
	"github.com/alebik0/go-auth/identity-service/data"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4/testutils/assert"

	_ "github.com/lib/pq" // To register the driver.

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
)

type ApiTest struct {
	Pmock         sqlmock.Sqlmock
	Rmock         redismock.ClientMock
	Cookies       []*http.Cookie
	Authorization string
	Router        *gin.Engine
}

func (apiTest *ApiTest) register(id uint32, login, password string, t *testing.T) {
	t.Logf("ApiTest.register")
	t.Logf("Prepare mock")
	apiTest.Pmock.
		ExpectQuery(
			regexp.QuoteMeta(
				"SELECT id, login, password_hash, profile_id FROM auth WHERE login = $1;",
			),
		).
		WithArgs(login).
		WillReturnError(sql.ErrNoRows)
	apiTest.Pmock.ExpectBegin()
	apiTest.Pmock.
		ExpectQuery(
			regexp.QuoteMeta(
				"INSERT INTO users (name, description) VALUES ($1, $2) RETURNING id;",
			),
		).
		WithArgs(login, "").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	apiTest.Pmock.
		ExpectQuery(
			regexp.QuoteMeta(
				"INSERT INTO auth (login, password_hash, profile_id) VALUES ($1, $2, $3) RETURNING id;",
			),
		).
		WithArgs(login, sqlmock.AnyArg(), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	apiTest.Pmock.ExpectCommit()

	t.Logf("Prepare redis mock")
	apiTest.Rmock.
		Regexp().
		ExpectSet(`^token:.+$`, strconv.FormatUint(uint64(id), 10), 30*24*time.Hour).
		SetVal("OK")

	t.Logf("POST /api/v1/auth/register")
	body := strings.NewReader(fmt.Sprintf(`{"login":"%s", "password":"%s"}`, login, password))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
	for _, c := range apiTest.Cookies {
		request.AddCookie(c)
	}
	request.Header.Add("Authorization", "Bearer "+apiTest.Authorization)
	responseWriter := httptest.NewRecorder()

	apiTest.Router.ServeHTTP(responseWriter, request)

	if responseWriter.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", responseWriter.Code)
	}

	var response api.RegisterResponse
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
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

	if err := apiTest.Pmock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	if err := apiTest.Rmock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	t.Logf("Update cookies")
	apiTest.Cookies = append(apiTest.Cookies, cookies...)
	apiTest.Authorization = response.AccessToken
}

func (apiTest *ApiTest) expect(id uint32, name, description string, expect data.UserData, t *testing.T) {
	t.Logf("ApiTest.expect")
	t.Logf("Prepare mock")
	apiTest.Pmock.
		ExpectQuery(
			regexp.QuoteMeta(
				"SELECT id, name, description FROM users WHERE id = $1;",
			),
		).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(id, name, description))

	t.Logf("GET /api/v1/users/my")
	body := http.NoBody
	request := httptest.NewRequest(http.MethodGet, "/api/v1/users/my", body)
	for _, c := range apiTest.Cookies {
		request.AddCookie(c)
	}
	request.Header.Add("Authorization", "Bearer "+apiTest.Authorization)

	responseWriter := httptest.NewRecorder()

	apiTest.Router.ServeHTTP(responseWriter, request)

	if responseWriter.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", responseWriter.Code, responseWriter.Body.Bytes())
	}

	var response data.UserData
	err := json.Unmarshal(responseWriter.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(response, expect) {
		t.Fatalf("got %+v want %+v", response, expect)
	}

	if err := apiTest.Pmock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	if err := apiTest.Rmock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	t.Logf("Update cookies")
	cookies := responseWriter.Result().Cookies()
	apiTest.Cookies = append(apiTest.Cookies, cookies...)
}

func (apiTest *ApiTest) expect401(t *testing.T) {
	t.Logf("ApiTest.expect401")

	t.Logf("GET /api/v1/users/my")
	body := http.NoBody
	request := httptest.NewRequest(http.MethodGet, "/api/v1/users/my", body)
	for _, c := range apiTest.Cookies {
		request.AddCookie(c)
	}
	request.Header.Add("Authorization", "Bearer "+apiTest.Authorization)

	responseWriter := httptest.NewRecorder()

	apiTest.Router.ServeHTTP(responseWriter, request)

	if responseWriter.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", responseWriter.Code, responseWriter.Body.Bytes())
	}

	if err := apiTest.Pmock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	if err := apiTest.Rmock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	t.Logf("Update cookies")
	cookies := responseWriter.Result().Cookies()
	apiTest.Cookies = append(apiTest.Cookies, cookies...)
}

func (apiTest *ApiTest) logout(t *testing.T) {
	t.Logf("ApiTest.logout")
	t.Logf("Prepare mock")

	t.Logf("Prepare redis mock")
	apiTest.Rmock.
		Regexp().
		ExpectDel(`^token:.+$`).
		SetVal(1)

	t.Logf("POST /api/v1/auth/logout")
	body := http.NoBody
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", body)
	for _, c := range apiTest.Cookies {
		request.AddCookie(c)
	}
	request.Header.Add("Authorization", "Bearer "+apiTest.Authorization)

	responseWriter := httptest.NewRecorder()

	apiTest.Router.ServeHTTP(responseWriter, request)

	if responseWriter.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", responseWriter.Code, responseWriter.Body.Bytes())
	}

	if err := apiTest.Pmock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	if err := apiTest.Rmock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}

	t.Logf("Update cookies")
	cookies := responseWriter.Result().Cookies()
	apiTest.Cookies = append(apiTest.Cookies, cookies...)
	apiTest.Authorization = ""
}

func NewApiTest(t *testing.T) (ApiTest, error) {
	gin.SetMode(gin.TestMode)

	t.Logf("Prepare data")
	hmacSecret := "qweqeqwioueqiewuio"

	t.Logf("Create postgres API mock")
	database, mock, err := sqlmock.New()
	if err != nil {
		return ApiTest{}, err
	}

	t.Logf("Create in-memory JWT API")
	redis, rmock := redismock.NewClientMock()

	t.Logf("Create handler")
	handler, err := api.NewHandler(
		database,
		redis,
		[]byte(hmacSecret),
	)
	if err != nil {
		return ApiTest{}, err
	}

	t.Logf("Setup router")
	router := SetupRouter(handler)

	return ApiTest{
		Pmock:         mock,
		Rmock:         rmock,
		Cookies:       make([]*http.Cookie, 0),
		Authorization: "",
		Router:        router,
	}, nil
}

func TestRegister(t *testing.T) {
	apiTest, err := NewApiTest(t)
	assert.NoError(t, err)

	login := "username"
	password := "passwd"

	apiTest.register(1, login, password, t)
	apiTest.expect(
		1,
		login,
		"",
		data.UserData{
			ID:          1,
			Name:        login,
			Description: "",
		},
		t,
	)
}

func TestLogout(t *testing.T) {
	apiTest, err := NewApiTest(t)
	assert.NoError(t, err)

	var id uint32 = 1
	login := "username"
	password := "passwd"

	apiTest.register(id, login, password, t)
	apiTest.expect(
		id,
		login,
		"",
		data.UserData{
			ID:          id,
			Name:        login,
			Description: "",
		},
		t,
	)
	apiTest.logout(t)
	apiTest.expect401(t)
}
