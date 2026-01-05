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

type PaymentApiControllerTestSuite struct {
	suite.Suite
	mockPaymentController *mockController.MockPaymentController
	apiController         *controller.PaymentApiController
	router                *chi.Mux
}

func (suite *PaymentApiControllerTestSuite) SetupTest() {
	suite.mockPaymentController = mockController.NewMockPaymentController(suite.T())
	suite.apiController = controller.NewPaymentApiController(suite.mockPaymentController)
	suite.router = chi.NewRouter()
	suite.apiController.RegisterRoutes(suite.router)
}

func TestPaymentApiControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentApiControllerTestSuite))
}

func (suite *PaymentApiControllerTestSuite) Test_CreatePayment_WithValidRequest_ShouldReturn201() {
	// GIVEN a valid payment request
	request := dto.AddPaymentRequestDto{
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
	}

	expectedQRCode := "00020101021243650016COM.MERCADOLIBRE"

	suite.mockPaymentController.EXPECT().
		CreatePayment(mock.Anything).
		Return(expectedQRCode, nil).
		Once()

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/v1/payment", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	// WHEN creating payment
	suite.router.ServeHTTP(rec, req)

	// THEN should return 201
	assert.Equal(suite.T(), http.StatusCreated, rec.Code)
	suite.mockPaymentController.AssertExpectations(suite.T())
}

func (suite *PaymentApiControllerTestSuite) Test_CreatePayment_WithInvalidJSON_ShouldReturn400() {
	// GIVEN invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/v1/payment", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	// WHEN creating payment with invalid JSON
	suite.router.ServeHTTP(rec, req)

	// THEN should return 400
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *PaymentApiControllerTestSuite) Test_CreatePayment_WithError_ShouldReturn500() {
	// GIVEN a payment request that will fail
	request := dto.AddPaymentRequestDto{
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
	}

	suite.mockPaymentController.EXPECT().
		CreatePayment(mock.Anything).
		Return("", errors.New("payment failed")).
		Once()

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/v1/payment", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	// WHEN payment creation fails
	suite.router.ServeHTTP(rec, req)

	// THEN should return 500
	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
	suite.mockPaymentController.AssertExpectations(suite.T())
}

func (suite *PaymentApiControllerTestSuite) Test_GetPaymentStatusByOrderId_WithValidId_ShouldReturn200() {
	// GIVEN a valid order ID
	orderId := uint(1)
	expectedStatus := "Approved"

	suite.mockPaymentController.EXPECT().
		GetPaymentStatusByOrderId(orderId).
		Return(expectedStatus, nil).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/v1/payment/1/status", nil)
	rec := httptest.NewRecorder()

	// WHEN getting payment status
	suite.router.ServeHTTP(rec, req)

	// THEN should return 200
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	suite.mockPaymentController.AssertExpectations(suite.T())
}

func (suite *PaymentApiControllerTestSuite) Test_GetPaymentStatusByOrderId_WithInvalidId_ShouldReturn400() {
	// GIVEN an invalid order ID
	req := httptest.NewRequest(http.MethodGet, "/v1/payment/invalid/status", nil)
	rec := httptest.NewRecorder()

	// WHEN getting payment status with invalid ID
	suite.router.ServeHTTP(rec, req)

	// THEN should return 400
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *PaymentApiControllerTestSuite) Test_GetPaymentStatusByOrderId_WithError_ShouldReturn500() {
	// GIVEN an order ID that will fail
	orderId := uint(999)

	suite.mockPaymentController.EXPECT().
		GetPaymentStatusByOrderId(orderId).
		Return("", errors.New("payment not found")).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/v1/payment/999/status", nil)
	rec := httptest.NewRecorder()

	// WHEN getting payment status fails
	suite.router.ServeHTTP(rec, req)

	// THEN should return 500
	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
	suite.mockPaymentController.AssertExpectations(suite.T())
}

func (suite *PaymentApiControllerTestSuite) Test_GetPaymentByOrderId_WithValidId_ShouldReturn200() {
	// GIVEN a valid order ID
	orderId := uint(1)
	expectedPayment := &dto.GetPaymentResponseDto{
		ID:      1,
		OrderId: orderId,
		Total:   100.50,
		Status:  "Approved",
		Type:    "QRCode",
	}

	suite.mockPaymentController.EXPECT().
		GetPaymentByOrderId(orderId).
		Return(expectedPayment, nil).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/v1/payment/1", nil)
	rec := httptest.NewRecorder()

	// WHEN getting payment
	suite.router.ServeHTTP(rec, req)

	// THEN should return 200
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	suite.mockPaymentController.AssertExpectations(suite.T())
}

func (suite *PaymentApiControllerTestSuite) Test_GetPaymentByOrderId_WithInvalidId_ShouldReturn400() {
	// GIVEN an invalid order ID
	req := httptest.NewRequest(http.MethodGet, "/v1/payment/invalid", nil)
	rec := httptest.NewRecorder()

	// WHEN getting payment with invalid ID
	suite.router.ServeHTTP(rec, req)

	// THEN should return 400
	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

func (suite *PaymentApiControllerTestSuite) Test_GetPaymentByOrderId_WithError_ShouldReturn500() {
	// GIVEN an order ID that will fail
	orderId := uint(999)

	suite.mockPaymentController.EXPECT().
		GetPaymentByOrderId(orderId).
		Return(nil, errors.New("payment not found")).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/v1/payment/999", nil)
	rec := httptest.NewRecorder()

	// WHEN getting payment fails
	suite.router.ServeHTTP(rec, req)

	// THEN should return 500
	assert.Equal(suite.T(), http.StatusInternalServerError, rec.Code)
	suite.mockPaymentController.AssertExpectations(suite.T())
}
