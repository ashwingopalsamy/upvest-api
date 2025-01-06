package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ashwingopalsamy/upvest-api/internal/domain"
)

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

func (r *userRepo) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, queryReadUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var (
			user          domain.User
			nationalities string
			postalAddress sql.NullString
			address       string
		)

		if err := rows.Scan(
			&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.FirstName, &user.LastName,
			&user.Salutation, &user.Title, &user.BirthDate, &user.BirthCity, &user.BirthCountry,
			&user.BirthName, &nationalities, &postalAddress, &address, &user.Status,
		); err != nil {
			return nil, err
		}

		// Deserializing the JSON fields
		if err := json.Unmarshal([]byte(nationalities), &user.Nationalities); err != nil {
			return nil, fmt.Errorf("failed to unmarshal nationalities: %w", err)
		}

		if postalAddress.Valid {
			if err := json.Unmarshal([]byte(postalAddress.String), &user.PostalAddress); err != nil {
				return nil, fmt.Errorf("failed to unmarshal postal_address: %w", err)
			}
		}

		if err := json.Unmarshal([]byte(address), &user.Address); err != nil {
			return nil, fmt.Errorf("failed to unmarshal address: %w", err)
		}

		users = append(users, user)
	}
	return users, rows.Err()
}
