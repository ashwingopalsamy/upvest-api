package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ashwingopalsamy/upvest-api/internal/domain"
	"github.com/ashwingopalsamy/upvest-api/internal/event"
	"github.com/ashwingopalsamy/upvest-api/internal/pkg/repository"
	"github.com/ashwingopalsamy/upvest-api/internal/util/writer"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type UserHandler struct {
	repo      repository.UserRepository
	publisher event.PublisherInterface
}

func NewUserHandler(repo repository.UserRepository, publisher event.PublisherInterface) *UserHandler {
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
	// Parse query parameters
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Default values
	offset := 0
	limit := 100

	// Parse integers with error handling
	if offsetVal, err := strconv.Atoi(offsetStr); err == nil && offsetVal >= 0 {
		offset = offsetVal
	}
	if limitVal, err := strconv.Atoi(limitStr); err == nil && limitVal > 0 && limitVal <= 1000 {
		limit = limitVal
	}

	users, err := h.repo.GetAllUsers(r.Context(), offset, limit, sort, order)
	if err != nil {
		writer.WriteErrJSON(w, http.StatusInternalServerError, ErrTitleDatabaseError, ErrMsgFailedToFetchUsers)
		return
	}

	writer.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"meta": map[string]interface{}{
			"count":  len(users),
			"offset": offset,
			"limit":  limit,
			"sort":   sort,
			"order":  order,
		},
		"data": users,
	})
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		writer.WriteErrJSON(w, http.StatusBadRequest, "invalid_request", "user_id is required")
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writer.WriteErrJSON(w, http.StatusNotFound, "not_found", "user does not exist")
		} else {
			writer.WriteErrJSON(w, http.StatusInternalServerError, "database_error", "failed to fetch user")
		}
		return
	}

	writer.WriteJSON(w, http.StatusOK, user)
}
