package main

import (
	"database/sql"
	"net/http"

	"github.com/ashwingopalsamy/upvest-api/internal/kafka"
	"github.com/ashwingopalsamy/upvest-api/internal/pkg/handler"
	"github.com/ashwingopalsamy/upvest-api/internal/pkg/repository"
	"github.com/gorilla/mux"
)

func NewServer(db *sql.DB, publisher *kafka.Publisher) http.Handler {
	router := mux.NewRouter()

	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo, *publisher)

	router.HandleFunc("/health", pingHTTP).Methods("GET")

	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")

	return router
}
