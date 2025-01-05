package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ashwingopalsamy/upvest-api/internal/domain"
	"github.com/ashwingopalsamy/upvest-api/internal/kafka"
	"github.com/ashwingopalsamy/upvest-api/internal/pkg/repository"
	"github.com/ashwingopalsamy/upvest-api/internal/pkg/util/writer"
	log "github.com/sirupsen/logrus"
)

type UserHandler struct {
	repo      repository.UserRepository
	publisher kafka.Publisher
}

func NewUserHandler(repo repository.UserRepository, publisher kafka.Publisher) *UserHandler {
	return &UserHandler{
		repo:      repo,
		publisher: publisher,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writer.WriteErrJSON(w, http.StatusBadRequest, "Invalid Request", "request body could not be parsed")
		return
	}

	if err := user.Validate(); err != nil {
		writer.WriteErrJSON(w, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	createdUser, err := h.repo.CreateUser(r.Context(), &user)
	if err != nil {
		log.Error(err)
		writer.WriteErrJSON(w, http.StatusInternalServerError, "Database Error", "failed to create user")
		return
	}

	// Prepare Kafka event and serialize it
	event := map[string]interface{}{
		"action": "USER_CREATED",
		"user":   createdUser,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		writer.WriteErrJSON(w, http.StatusInternalServerError, "Kafka Error", "failed to marshal event")
	}

	if err := h.publisher.Publish(r.Context(), []byte(createdUser.ID), eventBytes); err != nil {
		writer.WriteErrJSON(w, http.StatusInternalServerError, "Kafka Error", "failed to emit user creation event")
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdUser)
	if err != nil {
		return
	}
}
