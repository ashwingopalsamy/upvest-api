//go:generate mockery --name=UserRepository --output=../../util/mocks --outpkg=mocks
package repository

import (
	"context"
	"database/sql"

	"github.com/ashwingopalsamy/upvest-api/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}
