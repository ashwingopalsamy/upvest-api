package main

import (
	"database/sql"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// pingHTTP checks the health of our HTTP service.
func pingHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("200 Subscriber OK"))
	if err != nil {
		return
	}
}

// pingDB checks the health of our Postgres database.
func pingDB(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		log.Errorf("database health check failed: %v", err)
		return errors.New("database is unreachable")
	}
	return nil
}
