package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashwingopalsamy/upvest-api/internal/domain"
	"github.com/ashwingopalsamy/upvest-api/internal/pkg/handler"
	"github.com/ashwingopalsamy/upvest-api/internal/util/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	mockRepo      *mocks.UserRepository
	mockPublisher *mocks.PublisherInterface
	handler       *handler.UserHandler
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (suite *UserHandlerTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.UserRepository)
	suite.mockPublisher = new(mocks.PublisherInterface)
	suite.handler = handler.NewUserHandler(suite.mockRepo, suite.mockPublisher)
}

func (suite *UserHandlerTestSuite) TestCreateUser_Success() {
	ctx := context.Background()

	reqBody := &domain.User{
		FirstName:     "Rob",
		LastName:      "Schmidt",
		BirthDate:     "1990-01-01",
		BirthCity:     "Berlin",
		BirthCountry:  "DE",
		Nationalities: []string{"DE", "US"},
		PostalAddress: &domain.Address{
			AddressLine1: "123 Main St",
			Postcode:     "12345",
			City:         "Berlin",
			Country:      "DE",
		},
		Address: domain.Address{
			AddressLine1: "123 Main St",
			Postcode:     "12345",
			City:         "Berlin",
			Country:      "DE",
		},
	}

	suite.mockRepo.On("CreateUser", ctx, mock.Anything).Return(reqBody, nil)
	suite.mockPublisher.On("Publish", ctx, mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	w := httptest.NewRecorder()

	suite.handler.CreateUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	suite.Equal(http.StatusCreated, res.StatusCode)

	var resp domain.User
	err := json.NewDecoder(res.Body).Decode(&resp)
	suite.NoError(err)
	suite.Equal("Rob", resp.FirstName)

	// Assert expectations
	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockPublisher.AssertExpectations(suite.T())
}

func (suite *UserHandlerTestSuite) TestCreateUser_Failure() {
	reqBody := &domain.User{
		FirstName: "John",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	w := httptest.NewRecorder()

	suite.handler.CreateUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	suite.Equal(http.StatusBadRequest, res.StatusCode)
}
