//go:generate mockery --name=UserRepository --output=../../util/mocks --outpkg=mocks
package repository

import (
	"context"
	"database/sql"

	"github.com/ashwingopalsamy/upvest-api/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetAllUsers(ctx context.Context, offset, limit int, sort, order string) ([]domain.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}
