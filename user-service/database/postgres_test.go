package database

import (
	"regexp"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	_ "github.com/lib/pq" // To register the driver.

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	t.Logf("Create postgres API mock")
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	api := &PostgresAPI{database: db}
	defer func() {
		innerErr := api.Close()
		if innerErr != nil {
			t.Logf("[WARN] Failed to close database: %v", innerErr)
		}
	}()

	t.Logf("Prepare mock")
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.
		ExpectQuery(
			regexp.QuoteMeta(
				"INSERT INTO users (name, description) VALUES ($1, $2) RETURNING id;",
			),
		).
		WithArgs("hello", "world").
		WillReturnRows(rows)

	t.Logf("Create user")
	userData, err := api.CreateUser("hello", "world")
	expected := UserData{
		ID:          1,
		Name:        "hello",
		Description: "world",
	}

	t.Logf("Check result")
	require.NoError(t, err)
	require.Equal(t, expected, userData)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRead(t *testing.T) {
	t.Logf("Create postgres API mock")
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	api := &PostgresAPI{database: db}
	defer func() {
		innerErr := api.Close()
		if innerErr != nil {
			t.Logf("[WARN] Failed to close database: %v", innerErr)
		}
	}()

	t.Logf("Prepare mock")
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, "hello", "world")
	mock.
		ExpectQuery(
			regexp.QuoteMeta(
				"SELECT id, name, description FROM users WHERE id = $1;",
			),
		).
		WithArgs(1).
		WillReturnRows(rows)

	t.Logf("Read user")
	userData, err := api.ReadUser(1)
	expected := UserData{
		ID:          1,
		Name:        "hello",
		Description: "world",
	}

	t.Logf("Check result")
	require.NoError(t, err)
	require.Equal(t, expected, userData)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	t.Logf("Create postgres API mock")
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	api := &PostgresAPI{database: db}
	defer func() {
		innerErr := api.Close()
		if innerErr != nil {
			t.Logf("[WARN] Failed to close database: %v", innerErr)
		}
	}()

	t.Logf("Prepare mock")
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, "hello", "new world")
	mock.
		ExpectQuery(
			regexp.QuoteMeta(
				"UPDATE users SET name = $1, description = $2 WHERE id=$3 RETURNING id, name, description;",
			),
		).
		WithArgs("hello", "new world", 1).
		WillReturnRows(rows)

	t.Logf("Read user")
	userData, err := api.UpdateUser(1, "hello", "new world")
	expected := UserData{
		ID:          1,
		Name:        "hello",
		Description: "new world",
	}

	t.Logf("Check result")
	require.NoError(t, err)
	require.Equal(t, expected, userData)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	t.Logf("Create postgres API mock")
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	api := &PostgresAPI{database: db}
	defer func() {
		innerErr := api.Close()
		if innerErr != nil {
			t.Logf("[WARN] Failed to close database: %v", innerErr)
		}
	}()

	t.Logf("Prepare mock")
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, "hello", "world")
	mock.
		ExpectQuery(
			regexp.QuoteMeta(
				"DELETE FROM users WHERE id=$1 RETURNING id, name, description;",
			),
		).
		WithArgs(1).
		WillReturnRows(rows)

	t.Logf("Read user")
	userData, err := api.DeleteUser(1)
	expected := UserData{
		ID:          1,
		Name:        "hello",
		Description: "world",
	}

	t.Logf("Check result")
	require.NoError(t, err)
	require.Equal(t, expected, userData)

	require.NoError(t, mock.ExpectationsWereMet())
}
