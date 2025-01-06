package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

const publisherPortAddr = ":8080"

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
	log.Infof("starting server on %s", publisherPortAddr)
	if err := http.ListenAndServe(publisherPortAddr, server); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
