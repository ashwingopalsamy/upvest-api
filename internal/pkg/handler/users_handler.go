package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ashwingopalsamy/upvest-api/internal/domain"
	"github.com/ashwingopalsamy/upvest-api/internal/kafka"
	"github.com/ashwingopalsamy/upvest-api/internal/pkg/repository"
	"github.com/ashwingopalsamy/upvest-api/internal/util/writer"
	log "github.com/sirupsen/logrus"
)

type UserHandler struct {
	repo      repository.UserRepository
	publisher kafka.PublisherInterface
}

func NewUserHandler(repo repository.UserRepository, publisher kafka.PublisherInterface) *UserHandler {
	return &UserHandler{
		repo:      repo,
		publisher: publisher,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writer.WriteErrJSON(w, http.StatusBadRequest, ErrTitleInvalidRequest, ErrMsgInvalidRequestBody)
		return
	}

	if err := user.Validate(); err != nil {
		writer.WriteErrJSON(w, http.StatusBadRequest, ErrTitleValidationError, err.Error())
		return
	}

	createdUser, err := h.repo.CreateUser(r.Context(), &user)
	if err != nil {
		log.Error(err)
		writer.WriteErrJSON(w, http.StatusInternalServerError, ErrTitleDatabaseError, ErrMsgCreateUserFailed)
		return
	}

	// Prepare Kafka event and serialize it
	event := map[string]interface{}{
		"action": "USER_CREATED",
		"user":   createdUser,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		writer.WriteErrJSON(w, http.StatusInternalServerError, ErrTitleKafkaError, ErrMsgMarshalEventFailed)
		return
	}

	if err := h.publisher.Publish(r.Context(), []byte(createdUser.ID), eventBytes); err != nil {
		writer.WriteErrJSON(w, http.StatusInternalServerError, ErrTitleKafkaError, ErrMsgEmitEventFailed)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdUser)
	if err != nil {
		return
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.GetAllUsers(r.Context())
	if err != nil {
		writer.WriteErrJSON(w, http.StatusInternalServerError, ErrTitleDatabaseError, ErrMsgFailedToFetchUsers)
		return
	}
	writer.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"meta": map[string]interface{}{
			"count": len(users),
			"sort":  "created_at",
			"order": "ASC",
		},
		"data": users,
	})
}
