package repository

const fieldStatusActive = "ACTIVE"

var queryCreateUsers = `INSERT INTO users (first_name, last_name, salutation, title, birth_date, birth_city, birth_country,
                   birth_name, nationalities, postal_address, address, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::JSONB, $10::JSONB, $11::JSONB, $12)
RETURNING id, created_at, updated_at;`

var queryReadUsers = `SELECT id, created_at, updated_at, first_name, last_name, salutation, title, birth_date,
		       birth_city, birth_country, birth_name, nationalities, postal_address, address, status
FROM users
ORDER BY %s %s
LIMIT $1 OFFSET $2`

var queryOffboardUser = `UPDATE users 
		SET status = $1, updated_at = NOW()
		WHERE id = $2`
