//go:generate mockery --name=UserRepository --output=../../util/mocks --outpkg=mocks
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ashwingopalsamy/upvest-api/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetAllUsers(ctx context.Context, offset, limit int, sort, order string) ([]domain.User, error)
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	OffboardUser(ctx context.Context, userID string) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	postalAddress, _ := json.Marshal(user.PostalAddress)
	address, _ := json.Marshal(user.Address)
	nationalities, err := json.Marshal(user.Nationalities)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal nationalities: %w", err)
	}

	err = r.db.QueryRowContext(ctx, queryCreateUsers,
		user.FirstName, user.LastName, user.Salutation, user.Title,
		user.BirthDate, user.BirthCity, user.BirthCountry, user.BirthName,
		nationalities, postalAddress, address, fieldStatusActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) GetAllUsers(ctx context.Context, offset, limit int, sort, order string) ([]domain.User, error) {
	// Validate and normalize sorting inputs
	if sort != "created_at" && sort != "updated_at" {
		sort = "created_at"
	}
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	query := fmt.Sprintf(queryReadUsers, sort, order)

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var (
			user          domain.User
			nationalities sql.NullString
			postalAddress sql.NullString
			address       string
		)

		if err := rows.Scan(
			&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.FirstName, &user.LastName,
			&user.Salutation, &user.Title, &user.BirthDate, &user.BirthCity, &user.BirthCountry,
			&user.BirthName, &nationalities, &postalAddress, &address, &user.Status,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Deserializing the JSON fields
		if nationalities.Valid && nationalities.String != "" {
			if err := json.Unmarshal([]byte(nationalities.String), &user.Nationalities); err != nil {
				return nil, fmt.Errorf("failed to unmarshal nationalities: %w", err)
			}
		} else {
			user.Nationalities = nil
		}

		if postalAddress.Valid && postalAddress.String != "" {
			var addr domain.Address
			if err := json.Unmarshal([]byte(postalAddress.String), &addr); err != nil {
				return nil, fmt.Errorf("failed to unmarshal postal_address: %w", err)
			}
			user.PostalAddress = &addr
		} else {
			user.PostalAddress = nil
		}

		if err := json.Unmarshal([]byte(address), &user.Address); err != nil {
			return nil, fmt.Errorf("failed to unmarshal address: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `SELECT id, created_at, updated_at, first_name, last_name, salutation, title, birth_date,
		birth_city, birth_country, birth_name, nationalities, postal_address, address, status
		FROM users WHERE id = $1`

	var user domain.User
	var nationalities sql.NullString
	var postalAddress sql.NullString
	var address string

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.FirstName, &user.LastName,
		&user.Salutation, &user.Title, &user.BirthDate, &user.BirthCity, &user.BirthCountry,
		&user.BirthName, &nationalities, &postalAddress, &address, &user.Status,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Deserializing the JSON fields
	if nationalities.Valid {
		if err := json.Unmarshal([]byte(nationalities.String), &user.Nationalities); err != nil {
			return nil, fmt.Errorf("failed to unmarshal nationalities: %w", err)
		}
	}
	if postalAddress.Valid {
		if err := json.Unmarshal([]byte(postalAddress.String), &user.PostalAddress); err != nil {
			return nil, fmt.Errorf("failed to unmarshal postal_address: %w", err)
		}
	}
	if err := json.Unmarshal([]byte(address), &user.Address); err != nil {
		return nil, fmt.Errorf("failed to unmarshal address: %w", err)
	}

	return &user, nil
}

func (r *userRepo) OffboardUser(ctx context.Context, userID string) error {
	result, err := r.db.ExecContext(ctx, queryOffboardUser, "OFFBOARDED", userID)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
