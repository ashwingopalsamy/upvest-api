package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	DbDSN string
}

func main() {
	// Setup Logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.Info("starting Upvest API service")

	// Parse configuration
	config := Config{
		DbDSN: os.Getenv("DB_DSN"),
	}

	// Init Database
	db, err := initDatabase(config.DbDSN)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Init Kafka Publisher
	initKafkaPublisher()
	defer publisher.Close()

	// Create and start the HTTP server
	server := NewServer(db, publisher)

	// Init HTTP Server
	log.Info("starting server on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
