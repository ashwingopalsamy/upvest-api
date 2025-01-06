package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ashwingopalsamy/upvest-api/internal/domain"
	"github.com/stretchr/testify/assert"
)

var (
	db       *sql.DB
	mock     sqlmock.Sqlmock
	repo     UserRepository
	setupErr error
)

func setup() {
	db, mock, setupErr = sqlmock.New()
	if setupErr != nil {
		panic(setupErr)
	}
	repo = NewUserRepository(db)
}

func teardown() {
	db.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

// Test_CreateUser_Success
func Test_CreateUser_Success(t *testing.T) {
	mockUser := &domain.User{
		FirstName:     "Rob",
		LastName:      "Smith",
		BirthDate:     "1990-01-01",
		BirthCity:     "Berlin",
		BirthCountry:  "DE",
		Nationalities: []string{"DE", "US"},
		Address: domain.Address{
			AddressLine1: "123 Main St",
			Postcode:     "12345",
			City:         "Berlin",
			Country:      "DE",
		},
	}

	nationalities, _ := json.Marshal(mockUser.Nationalities)
	address, _ := json.Marshal(mockUser.Address)

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(
			mockUser.FirstName,
			mockUser.LastName,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			mockUser.BirthDate,
			mockUser.BirthCity,
			mockUser.BirthCountry,
			sqlmock.AnyArg(),
			nationalities,
			sqlmock.AnyArg(),
			address,
			"ACTIVE",
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("123", "2025-01-01T00:00:00Z", "2025-01-01T00:00:00Z"))

	createdUser, err := repo.CreateUser(context.Background(), mockUser)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, "123", createdUser.ID)
	assert.Equal(t, "2025-01-01T00:00:00Z", createdUser.CreatedAt)
}

// Test_CreateUser_Failure tests the failure case for CreateUser
func Test_CreateUser_Failure(t *testing.T) {
	mockUser := &domain.User{
		FirstName:     "Rob",
		LastName:      "Smith",
		BirthDate:     "1999-01-01",
		BirthCity:     "Berlin",
		BirthCountry:  "DE",
		Nationalities: []string{"DE", "US"},
		Address: domain.Address{
			AddressLine1: "123 Main St",
			Postcode:     "12345",
			City:         "Berlin",
			Country:      "DE",
		},
	}

	mock.ExpectQuery(`INSERT INTO users`).WillReturnError(sql.ErrConnDone)

	createdUser, err := repo.CreateUser(context.Background(), mockUser)

	assert.Error(t, err)
	assert.Nil(t, createdUser)
	assert.EqualError(t, err, sql.ErrConnDone.Error())
}

// Test_GetAllUsers_Success tests the success case for GetAllUsers
func Test_GetAllUsers_Success(t *testing.T) {
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "first_name", "last_name", "salutation", "title", "birth_date",
		"birth_city", "birth_country", "birth_name", "nationalities", "postal_address", "address", "status",
	}).
		AddRow("1", "2025-01-01T00:00:00Z", "2025-01-01T00:00:00Z", "John", "Schmidt", "", "DR", "1998-01-01",
			"Berlin", "DE", "", `["DE"]`, `{"address_line1":"123 Main St"}`, `{"address_line1":"456 High St"}`, "ACTIVE").
		AddRow("2", "2025-01-02T00:00:00Z", "2025-01-02T00:00:00Z", "Jane", "Schmidt", "", "PROF", "1999-01-01",
			"Munich", "DE", "", `["DE","US"]`, `{"address_line1":"789 Park Ave"}`, `{"address_line1":"123 High St"}`, "ACTIVE")

	mock.ExpectQuery(`SELECT id, created_at, updated_at`).
		WithArgs(100, 0).
		WillReturnRows(rows)

	users, err := repo.GetAllUsers(context.Background(), 0, 100, "created_at", "ASC")

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, "1", users[0].ID)
	assert.Equal(t, "Jane", users[1].FirstName)
	assert.Equal(t, []string{"DE", "US"}, users[1].Nationalities)
}

// Test_GetAllUsers_Failure tests the failure case for GetAllUsers
func Test_GetAllUsers_Failure(t *testing.T) {
	mock.ExpectQuery(`SELECT id, created_at, updated_at`).WillReturnError(sql.ErrConnDone)

	users, err := repo.GetAllUsers(context.Background(), 0, 100, "created_at", "ASC")

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.EqualError(t, err, "failed to execute query: sql: connection is already closed")
}

// Test_GetAllUsers_InvalidSorting tests invalid sorting and ensures it defaults to created_at
func Test_GetAllUsers_InvalidSorting(t *testing.T) {
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "first_name", "last_name", "salutation", "title", "birth_date",
		"birth_city", "birth_country", "birth_name", "nationalities", "postal_address", "address", "status",
	}).AddRow("1", "2025-01-01T00:00:00Z", "2025-01-01T00:00:00Z", "Jason", "Schmidt", "", "", "2000-01-01",
		"Berlin", "DE", "", `["DE"]`, `{"address_line1":"123 Main St"}`, `{"address_line1":"456 High St"}`, "ACTIVE")

	mock.ExpectQuery(`SELECT id, created_at, updated_at, first_name, last_name, salutation, title, birth_date, 
		       birth_city, birth_country, birth_name, nationalities, postal_address, address, status 
		FROM users ORDER BY created_at ASC LIMIT \$1 OFFSET \$2`).
		WithArgs(100, 0).
		WillReturnRows(rows)

	users, err := repo.GetAllUsers(context.Background(), 0, 100, "invalid_field", "invalid_order")

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 1)
	assert.Equal(t, "Jason", users[0].FirstName)
}

// Test_GetAllUsers_InvalidPagination tests invalid pagination parameters
func Test_GetAllUsers_InvalidPagination(t *testing.T) {
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "first_name", "last_name", "salutation", "title", "birth_date",
		"birth_city", "birth_country", "birth_name", "nationalities", "postal_address", "address", "status",
	}).AddRow("1", "2025-01-01T00:00:00Z", "2025-01-01T00:00:00Z", "Mark", "Smith", "", "", "1985-01-01",
		"Berlin", "DE", "", `["DE"]`, `{"address_line1":"789 Main St"}`, `{"address_line1":"123 Side St"}`, "ACTIVE")

	mock.ExpectQuery(`SELECT id, created_at, updated_at, first_name, last_name, salutation, title, birth_date, 
		       birth_city, birth_country, birth_name, nationalities, postal_address, address, status 
		FROM users ORDER BY created_at ASC LIMIT \$1 OFFSET \$2`).
		WithArgs(200, 0).
		WillReturnRows(rows)

	users, err := repo.GetAllUsers(context.Background(), 0, 200, "created_at", "ASC")

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 1)
	assert.Equal(t, "Mark", users[0].FirstName)
}

func Test_GetUserByID_Success(t *testing.T) {
	row := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "first_name", "last_name", "salutation", "title", "birth_date",
		"birth_city", "birth_country", "birth_name", "nationalities", "postal_address", "address", "status",
	}).AddRow("1", "2025-01-01T00:00:00Z", "2025-01-01T00:00:00Z", "Jason", "Schmidt", "SALUTATION_MALE", "DR",
		"2001-01-01", "Berlin", "DE", "Schmidt", `["DE"]`, `{"address_line1":"123 Main St"}`,
		`{"address_line1":"123 Main St"}`, "ACTIVE")

	mock.ExpectQuery(`SELECT id, created_at, updated_at`).WithArgs("1").WillReturnRows(row)

	user, err := repo.GetUserByID(context.Background(), "1")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Jason", user.FirstName)
	assert.Equal(t, "Schmidt", user.LastName)
}

func Test_GetUserByID_NotFound(t *testing.T) {
	mock.ExpectQuery(`SELECT id, created_at, updated_at`).WithArgs("1").WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByID(context.Background(), "1")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}

func Test_OffboardUser_Success(t *testing.T) {
	setup()
	defer teardown()

	mock.ExpectExec(`UPDATE users SET status = \$1, updated_at = NOW\(\) WHERE id = \$2`).
		WithArgs("OFFBOARDED", "123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.OffboardUser(context.Background(), "123")

	assert.NoError(t, err)
}

func Test_OffboardUser_NotFound(t *testing.T) {
	setup()
	defer teardown()

	mock.ExpectExec(`UPDATE users SET status = \$1, updated_at = NOW\(\) WHERE id = \$2`).
		WithArgs("OFFBOARDED", "123").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.OffboardUser(context.Background(), "123")

	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}

func Test_OffboardUser_DatabaseError(t *testing.T) {
	setup()
	defer teardown()

	mock.ExpectExec(`UPDATE users SET status = \$1, updated_at = NOW\(\) WHERE id = \$2`).
		WithArgs("OFFBOARDED", "123").
		WillReturnError(fmt.Errorf("database error"))

	err := repo.OffboardUser(context.Background(), "123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user status: database error")
}

