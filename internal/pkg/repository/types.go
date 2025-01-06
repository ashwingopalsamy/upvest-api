package repository

const fieldStatusActive = "ACTIVE"

var queryInsertUsers = `INSERT INTO users (first_name, last_name, salutation, title, birth_date, birth_city, birth_country,
                   birth_name, nationalities, postal_address, address, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::JSONB, $10::JSONB, $11::JSONB, $12)
RETURNING id, created_at, updated_at;`
