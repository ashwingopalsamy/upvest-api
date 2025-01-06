package repository

import (
	"context"
	"database/sql"
	"encoding/json"
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
