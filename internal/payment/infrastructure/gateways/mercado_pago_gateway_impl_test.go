package gateways_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/gateways"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockHTTPClient is a mock for rest.HTTPClient
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

type MercadoPagoGatewayTestSuite struct {
	suite.Suite
	mockHTTPClient *MockHTTPClient
}

func (suite *MercadoPagoGatewayTestSuite) SetupTest() {
	suite.mockHTTPClient = new(MockHTTPClient)
	os.Setenv("MERCADO_PAGO_BASEURL", "https://api.mercadopago.com")
	os.Setenv("MERCADO_PAGO_ACCESS_TOKEN", "test_token")
	os.Setenv("MERCADO_PAGO_CLIENT_ID", "client123")
	os.Setenv("MERCADO_PAGO_POS_ID", "pos123")
}

func (suite *MercadoPagoGatewayTestSuite) TearDownTest() {
	os.Unsetenv("MERCADO_PAGO_BASEURL")
	os.Unsetenv("MERCADO_PAGO_ACCESS_TOKEN")
	os.Unsetenv("MERCADO_PAGO_CLIENT_ID")
	os.Unsetenv("MERCADO_PAGO_POS_ID")
}

func TestMercadoPagoGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(MercadoPagoGatewayTestSuite))
}

func (suite *MercadoPagoGatewayTestSuite) Test_GenerateQRCode_WithValidRequest_ShouldReturnQRCode() {
	// GIVEN a valid QR code request
	request := dto.CreateQRCodeDTO{
		ExternalReference: "order-1",
		Title:             "Test Order",
		Description:       "Test Description",
		NotificationURL:   "http://localhost:8082/webhook",
		TotalAmount:       100.50,
		Items: []dto.Item{
			{
				SKUNumber:   "1",
				Category:    "food",
				Title:       "Product",
				Description: "Description",
				UnitPrice:   50.25,
				Quantity:    2,
				TotalAmount: 100.50,
			},
		},
	}

	expectedResponse := dto.QRCodeResponseDto{
		QRData: "00020101021243650016COM.MERCADOLIBRE",
	}

	responseBody, _ := json.Marshal(expectedResponse)
	response := &http.Response{
		StatusCode: http.StatusCreated,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}

	suite.mockHTTPClient.On("Do", mock.Anything).Return(response, nil).Once()

	gateway, err := gateways.NewMercadoPagoGatewayImplWithClient(suite.mockHTTPClient)
	assert.NoError(suite.T(), err)

	// WHEN generating QR code
	result, err := gateway.GenerateQRCode(context.Background(), request)

	// THEN QR code should be generated
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedResponse.QRData, result.QRData)
}

func (suite *MercadoPagoGatewayTestSuite) Test_GenerateQRCode_WithHTTPError_ShouldReturnError() {
	// GIVEN a QR code request
	request := dto.CreateQRCodeDTO{
		ExternalReference: "order-1",
		TotalAmount:       100.50,
	}

	expectedError := errors.New("connection failed")
	suite.mockHTTPClient.On("Do", mock.Anything).Return(nil, expectedError).Once()

	gateway, err := gateways.NewMercadoPagoGatewayImplWithClient(suite.mockHTTPClient)
	assert.NoError(suite.T(), err)

	// WHEN HTTP request fails
	result, err := gateway.GenerateQRCode(context.Background(), request)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "HTTP request failed")
	assert.Empty(suite.T(), result.QRData)
}

func (suite *MercadoPagoGatewayTestSuite) Test_GenerateQRCode_WithNonCreatedStatus_ShouldReturnError() {
	// GIVEN a QR code request
	request := dto.CreateQRCodeDTO{
		ExternalReference: "order-1",
		TotalAmount:       100.50,
	}

	response := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewReader([]byte("invalid request"))),
	}

	suite.mockHTTPClient.On("Do", mock.Anything).Return(response, nil).Once()

	gateway, err := gateways.NewMercadoPagoGatewayImplWithClient(suite.mockHTTPClient)
	assert.NoError(suite.T(), err)

	// WHEN response status is not 201
	result, err := gateway.GenerateQRCode(context.Background(), request)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to generate QR code")
	assert.Empty(suite.T(), result.QRData)
}

func (suite *MercadoPagoGatewayTestSuite) Test_GenerateQRCode_WithInvalidJSON_ShouldReturnError() {
	// GIVEN a QR code request
	request := dto.CreateQRCodeDTO{
		ExternalReference: "order-1",
		TotalAmount:       100.50,
	}

	response := &http.Response{
		StatusCode: http.StatusCreated,
		Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
	}

	suite.mockHTTPClient.On("Do", mock.Anything).Return(response, nil).Once()

	gateway, err := gateways.NewMercadoPagoGatewayImplWithClient(suite.mockHTTPClient)
	assert.NoError(suite.T(), err)

	// WHEN response has invalid JSON
	result, err := gateway.GenerateQRCode(context.Background(), request)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to decode response body")
	assert.Empty(suite.T(), result.QRData)
}

func (suite *MercadoPagoGatewayTestSuite) Test_NewMercadoPagoGatewayImpl_WithMissingConfig_ShouldReturnError() {
	// GIVEN missing environment variables
	os.Unsetenv("MERCADO_PAGO_ACCESS_TOKEN")

	// WHEN creating gateway
	gateway, err := gateways.NewMercadoPagoGatewayImpl()

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), gateway)
	assert.Contains(suite.T(), err.Error(), "invalid MercadoPagoConfig")
}

func (suite *MercadoPagoGatewayTestSuite) Test_NewMercadoPagoGatewayImpl_WithValidConfig_ShouldSucceed() {
	// GIVEN valid environment variables (already set in SetupTest)

	// WHEN creating gateway
	gateway, err := gateways.NewMercadoPagoGatewayImpl()

	// THEN gateway should be created
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), gateway)
}
