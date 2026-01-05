package clients_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/clients"
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

type OrderClientTestSuite struct {
	suite.Suite
	mockHTTPClient *MockHTTPClient
	client         clients.OrderClient
}

func (suite *OrderClientTestSuite) SetupTest() {
	suite.mockHTTPClient = new(MockHTTPClient)
	os.Setenv("ORDER_SERVICE_URL", "http://localhost:8081")
	suite.client = clients.NewOrderClient(suite.mockHTTPClient)
}

func (suite *OrderClientTestSuite) TearDownTest() {
	os.Unsetenv("ORDER_SERVICE_URL")
}

func TestOrderClientTestSuite(t *testing.T) {
	suite.Run(t, new(OrderClientTestSuite))
}

func (suite *OrderClientTestSuite) Test_GetOrder_WithValidOrderId_ShouldReturnOrder() {
	// GIVEN a valid order ID
	orderId := uint(1)
	expectedOrder := &dto.OrderResponseDto{
		ID:          orderId,
		TotalAmount: 100.50,
		Products: []*dto.OrderProductDto{
			{
				ProductId: 1,
				Name:      "Product Test",
				Price:     50.25,
				Quantity:  2,
			},
		},
	}

	responseBody, _ := json.Marshal(expectedOrder)
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}

	suite.mockHTTPClient.On("Do", mock.Anything).Return(response, nil).Once()

	// WHEN getting order
	order, err := suite.client.GetOrder(orderId)

	// THEN order should be returned
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), order)
	assert.Equal(suite.T(), orderId, order.ID)
	assert.Equal(suite.T(), float32(100.50), order.TotalAmount)
	suite.mockHTTPClient.AssertExpectations(suite.T())
}

func (suite *OrderClientTestSuite) Test_GetOrder_WithHTTPError_ShouldReturnError() {
	// GIVEN an order ID
	orderId := uint(1)
	expectedError := errors.New("connection failed")

	suite.mockHTTPClient.On("Do", mock.Anything).Return(nil, expectedError).Once()

	// WHEN HTTP request fails
	order, err := suite.client.GetOrder(orderId)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), order)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockHTTPClient.AssertExpectations(suite.T())
}

func (suite *OrderClientTestSuite) Test_GetOrder_WithNonOKStatus_ShouldReturnError() {
	// GIVEN an order ID
	orderId := uint(999)

	response := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewReader([]byte("order not found"))),
	}

	suite.mockHTTPClient.On("Do", mock.Anything).Return(response, nil).Once()

	// WHEN order is not found
	order, err := suite.client.GetOrder(orderId)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), order)
	assert.Contains(suite.T(), err.Error(), "failed to get order")
	suite.mockHTTPClient.AssertExpectations(suite.T())
}

func (suite *OrderClientTestSuite) Test_GetOrder_WithInvalidJSON_ShouldReturnError() {
	// GIVEN an order ID
	orderId := uint(1)

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
	}

	suite.mockHTTPClient.On("Do", mock.Anything).Return(response, nil).Once()

	// WHEN response has invalid JSON
	order, err := suite.client.GetOrder(orderId)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), order)
	suite.mockHTTPClient.AssertExpectations(suite.T())
}

func (suite *OrderClientTestSuite) Test_UpdateOrderStatus_WithValidData_ShouldSucceed() {
	// GIVEN valid order ID and status
	orderId := uint(1)
	status := 2

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(""))),
	}

	suite.mockHTTPClient.On("Do", mock.Anything).Return(response, nil).Once()

	// WHEN updating order status
	err := suite.client.UpdateOrderStatus(orderId, status)

	// THEN operation should succeed
	assert.NoError(suite.T(), err)
	suite.mockHTTPClient.AssertExpectations(suite.T())
}

func (suite *OrderClientTestSuite) Test_UpdateOrderStatus_WithNoContentStatus_ShouldSucceed() {
	// GIVEN valid order ID and status
	orderId := uint(1)
	status := 2

	response := &http.Response{
		StatusCode: http.StatusNoContent,
		Body:       io.NopCloser(bytes.NewReader([]byte(""))),
	}

	suite.mockHTTPClient.On("Do", mock.Anything).Return(response, nil).Once()

	// WHEN updating order status
	err := suite.client.UpdateOrderStatus(orderId, status)

	// THEN operation should succeed
	assert.NoError(suite.T(), err)
	suite.mockHTTPClient.AssertExpectations(suite.T())
}

func (suite *OrderClientTestSuite) Test_UpdateOrderStatus_WithHTTPError_ShouldReturnError() {
	// GIVEN an order ID and status
	orderId := uint(1)
	status := 2
	expectedError := errors.New("connection failed")

	suite.mockHTTPClient.On("Do", mock.Anything).Return(nil, expectedError).Once()

	// WHEN HTTP request fails
	err := suite.client.UpdateOrderStatus(orderId, status)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockHTTPClient.AssertExpectations(suite.T())
}

func (suite *OrderClientTestSuite) Test_UpdateOrderStatus_WithNonOKStatus_ShouldReturnError() {
	// GIVEN an order ID and status
	orderId := uint(999)
	status := 2

	response := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewReader([]byte("update failed"))),
	}

	suite.mockHTTPClient.On("Do", mock.Anything).Return(response, nil).Once()

	// WHEN update fails
	err := suite.client.UpdateOrderStatus(orderId, status)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to update order status")
	suite.mockHTTPClient.AssertExpectations(suite.T())
}

func (suite *OrderClientTestSuite) Test_NewOrderClient_WithoutEnvVar_ShouldUseDefaultURL() {
	// GIVEN no ORDER_SERVICE_URL env var
	os.Unsetenv("ORDER_SERVICE_URL")

	// WHEN creating client
	client := clients.NewOrderClient(suite.mockHTTPClient)

	// THEN client should be created with default URL
	assert.NotNil(suite.T(), client)
}
