package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ashwingopalsamy/upvest-api/internal/domain"
)

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	postalAddress, _ := json.Marshal(user.PostalAddress)
	address, _ := json.Marshal(user.Address)
	nationalities, err := json.Marshal(user.Nationalities)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal nationalities: %w", err)
	}

	err = r.db.QueryRowContext(ctx, queryInsertUsers,
		user.FirstName, user.LastName, user.Salutation, user.Title,
		user.BirthDate, user.BirthCity, user.BirthCountry, user.BirthName,
		nationalities, postalAddress, address, "ACTIVE",
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}
