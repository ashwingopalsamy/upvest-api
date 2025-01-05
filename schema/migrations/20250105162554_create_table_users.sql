-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   first_name VARCHAR(100) NOT NULL,
   last_name VARCHAR(100) NOT NULL,
   salutation VARCHAR(50) CHECK (salutation IN ('', 'SALUTATION_MALE', 'SALUTATION_FEMALE', 'SALUTATION_FEMALE_MARRIED', 'SALUTATION_DIVERSE')) DEFAULT '',
   title VARCHAR(50) CHECK (title IN ('', 'DR', 'PROF', 'PROF_DR', 'DIPL_ING', 'MAGISTER')) DEFAULT '',
   birth_date DATE NOT NULL,
   birth_city VARCHAR(85),
   birth_country CHAR(2) NOT NULL,
   birth_name VARCHAR(100),
   nationalities JSONB NOT NULL, -- Changed from TEXT[] to JSONB
   postal_address JSONB,
   address JSONB NOT NULL,
   status VARCHAR(20) CHECK (status IN ('ACTIVE', 'INACTIVE', 'OFFBOARDING', 'OFFBOARDED')) DEFAULT 'ACTIVE'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

