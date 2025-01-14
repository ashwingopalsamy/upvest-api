package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const subscriberPortAddr = ":8081"

type Config struct {
	DbDSN string
}

func main() {
	// Setup Logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.Info("starting Upvest API Subscriber service")

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

	// Init Subscriber
	initKafkaSubscriber()
	defer subscriber.Close()

	go func() {
		subscriber.Consume(context.TODO(), kafkaListener)
	}()

	// Setup Router
	router := mux.NewRouter()

	// Setup Routes
	router.HandleFunc("/health", pingHTTP).Methods("GET")

	// Init HTTP Server
	log.Info("starting server on %d", subscriberPortAddr)
	if err := http.ListenAndServe(subscriberPortAddr, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
