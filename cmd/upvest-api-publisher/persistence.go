package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// initDatabase initializes and returns a database connection.
func initDatabase(dbURL string) (*sql.DB, error) {
	// Open the database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Verify the connection
	if err := pingDB(db); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Info("database connection established")
	return db, nil
}
