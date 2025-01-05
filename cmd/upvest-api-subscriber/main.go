package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	DBURL string
}

func main() {
	// Setup Logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.Info("starting Upvest API service")

	// Parse configuration
	config := Config{
		DBURL: os.Getenv("DBURL"),
	}

	// Init Database
	db, err := initDatabase(config.DBURL)
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
	log.Info("starting server on :8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
