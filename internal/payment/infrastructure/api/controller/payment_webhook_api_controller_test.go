package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/controller"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	mockController "github.com/abattassini/tc-fiap-payment/mocks/payment/controller"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PaymentWebhookApiControllerTestSuite struct {
	suite.Suite
	mockWebhookController *mockController.MockPaymentWebhookController
	apiController         *controller.PaymentWebhookApiController
	router                *chi.Mux
}

func (suite *PaymentWebhookApiControllerTestSuite) SetupTest() {
	suite.mockWebhookController = mockController.NewMockPaymentWebhookController(suite.T())
	suite.apiController = controller.NewPaymentWebhookApiController(suite.mockWebhookController)
	suite.router = chi.NewRouter()
	suite.apiController.RegisterRoutes(suite.router)
}

func TestPaymentWebhookApiControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentWebhookApiControllerTestSuite))
}

func (suite *PaymentWebhookApiControllerTestSuite) Test_HandlePaymentNotification_WithValidRequest_ShouldReturn200() {
	// GIVEN a valid webhook notification
	request := dto.MercadoPagoWebhookNotificationRequestDTO{
		Id:    "1",
		Topic: "payment.created",
	}

	suite.mockWebhookController.EXPECT().
		HandleWebhook(mock.MatchedBy(func(req *dto.MercadoPagoWebhookNotificationRequestDTO) bool {
			return req.Id == "1" && req.Topic == "payment.created"
		})).
		Return(nil).
		Once()

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/payment/webhooks/notify", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	// WHEN handling webhook
	suite.router.ServeHTTP(rec, req)

	// THEN should return 200
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	suite.mockWebhookController.AssertExpectations(suite.T())
}

func (suite *PaymentWebhookApiControllerTestSuite) Test_HandlePaymentNotification_WithInvalidJSON_ShouldReturn400() {
	// GIVEN invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/payment/webhooks/notify", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	// WHEN handling webhook with invalid JSON
	suite.router.ServeHTTP(rec, req)

	// THEN should return 400
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *PaymentWebhookApiControllerTestSuite) Test_HandlePaymentNotification_WithError_ShouldReturn500() {
	// GIVEN a webhook request that will fail
	request := dto.MercadoPagoWebhookNotificationRequestDTO{
		Id:    "1",
		Topic: "payment.created",
	}

	suite.mockWebhookController.EXPECT().
		HandleWebhook(mock.MatchedBy(func(req *dto.MercadoPagoWebhookNotificationRequestDTO) bool {
			return req.Id == "1" && req.Topic == "payment.created"
		})).
		Return(errors.New("webhook processing failed")).
		Once()

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/payment/webhooks/notify", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	// WHEN webhook processing fails
	suite.router.ServeHTTP(rec, req)

	// THEN should return 500
	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
	suite.mockWebhookController.AssertExpectations(suite.T())
}
