package controller_test

import (
	"errors"
	"testing"

	"github.com/abattassini/tc-fiap-payment/internal/payment/controller"
	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	mockPresenter "github.com/abattassini/tc-fiap-payment/mocks/payment/presenter"
	mockAddPayment "github.com/abattassini/tc-fiap-payment/mocks/payment/usecase/addPayment"
	mockGetPayment "github.com/abattassini/tc-fiap-payment/mocks/payment/usecase/getPayment"
	mockGetPaymentStatus "github.com/abattassini/tc-fiap-payment/mocks/payment/usecase/getPaymentStatus"
	mockUpdatePayment "github.com/abattassini/tc-fiap-payment/mocks/payment/usecase/updatePayment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PaymentControllerTestSuite struct {
	suite.Suite
	mockPresenter               *mockPresenter.MockPaymentPresenter
	mockGetPaymentUseCase       *mockGetPayment.MockGetPaymentUseCase
	mockGetPaymentStatusUseCase *mockGetPaymentStatus.MockGetPaymentStatusUseCase
	mockUpdatePaymentUseCase    *mockUpdatePayment.MockUpdatePaymentUseCase
	mockAddPaymentUseCase       *mockAddPayment.MockAddPaymentUseCase
	controller                  controller.PaymentController
}

func (suite *PaymentControllerTestSuite) SetupTest() {
	suite.mockPresenter = mockPresenter.NewMockPaymentPresenter(suite.T())
	suite.mockGetPaymentUseCase = mockGetPayment.NewMockGetPaymentUseCase(suite.T())
	suite.mockGetPaymentStatusUseCase = mockGetPaymentStatus.NewMockGetPaymentStatusUseCase(suite.T())
	suite.mockUpdatePaymentUseCase = mockUpdatePayment.NewMockUpdatePaymentUseCase(suite.T())
	suite.mockAddPaymentUseCase = mockAddPayment.NewMockAddPaymentUseCase(suite.T())
	suite.controller = controller.NewPaymentControllerImpl(
		suite.mockPresenter,
		suite.mockGetPaymentUseCase,
		suite.mockGetPaymentStatusUseCase,
		suite.mockUpdatePaymentUseCase,
		suite.mockAddPaymentUseCase,
	)
}

func TestPaymentControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentControllerTestSuite))
}

func (suite *PaymentControllerTestSuite) Test_CreatePayment_WithValidRequest_ShouldReturnQRCode() {
	// GIVEN a valid payment request
	request := &dto.AddPaymentRequestDto{
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
	}

	expectedQRCode := "00020101021243650016COM.MERCADOLIBRE"

	suite.mockAddPaymentUseCase.EXPECT().
		Execute(mock.MatchedBy(func(cmd interface{}) bool {
			return true
		})).
		Return(expectedQRCode, nil).
		Once()

	// WHEN creating payment
	qrCode, err := suite.controller.CreatePayment(request)

	// THEN QR code should be returned
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedQRCode, qrCode)
	suite.mockAddPaymentUseCase.AssertExpectations(suite.T())
}

func (suite *PaymentControllerTestSuite) Test_CreatePayment_WithError_ShouldReturnError() {
	// GIVEN a payment request
	request := &dto.AddPaymentRequestDto{
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
	}

	expectedError := errors.New("payment creation failed")

	suite.mockAddPaymentUseCase.EXPECT().
		Execute(mock.Anything).
		Return("", expectedError).
		Once()

	// WHEN payment creation fails
	qrCode, err := suite.controller.CreatePayment(request)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	assert.Empty(suite.T(), qrCode)
	suite.mockAddPaymentUseCase.AssertExpectations(suite.T())
}

func (suite *PaymentControllerTestSuite) Test_GetPaymentStatusByOrderId_WithValidOrderId_ShouldReturnStatus() {
	// GIVEN a valid order ID
	orderId := uint(1)
	expectedStatus := "Approved"

	suite.mockGetPaymentStatusUseCase.EXPECT().
		Execute(mock.Anything).
		Return(expectedStatus, nil).
		Once()

	// WHEN getting payment status
	status, err := suite.controller.GetPaymentStatusByOrderId(orderId)

	// THEN status should be returned
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedStatus, status)
	suite.mockGetPaymentStatusUseCase.AssertExpectations(suite.T())
}

func (suite *PaymentControllerTestSuite) Test_GetPaymentStatusByOrderId_WithError_ShouldReturnError() {
	// GIVEN an order ID
	orderId := uint(999)
	expectedError := errors.New("payment not found")

	suite.mockGetPaymentStatusUseCase.EXPECT().
		Execute(mock.Anything).
		Return("", expectedError).
		Once()

	// WHEN getting payment status fails
	status, err := suite.controller.GetPaymentStatusByOrderId(orderId)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	assert.Empty(suite.T(), status)
	suite.mockGetPaymentStatusUseCase.AssertExpectations(suite.T())
}

func (suite *PaymentControllerTestSuite) Test_GetPaymentByOrderId_WithValidOrderId_ShouldReturnPayment() {
	// GIVEN a valid order ID
	orderId := uint(1)
	payment := &entities.Payment{
		ID:      1,
		OrderId: orderId,
		Total:   100.50,
		Status:  "Approved",
		Type:    "QRCode",
	}

	expectedResponse := &dto.GetPaymentResponseDto{
		ID:      1,
		OrderId: orderId,
		Total:   100.50,
		Status:  "Approved",
		Type:    "QRCode",
	}

	suite.mockGetPaymentUseCase.EXPECT().
		Execute(mock.Anything).
		Return(payment, nil).
		Once()

	suite.mockPresenter.EXPECT().
		Present(payment).
		Return(expectedResponse).
		Once()

	// WHEN getting payment
	response, err := suite.controller.GetPaymentByOrderId(orderId)

	// THEN payment should be returned
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedResponse, response)
	suite.mockGetPaymentUseCase.AssertExpectations(suite.T())
	suite.mockPresenter.AssertExpectations(suite.T())
}

func (suite *PaymentControllerTestSuite) Test_GetPaymentByOrderId_WithError_ShouldReturnError() {
	// GIVEN an order ID
	orderId := uint(999)
	expectedError := errors.New("payment not found")

	suite.mockGetPaymentUseCase.EXPECT().
		Execute(mock.Anything).
		Return(nil, expectedError).
		Once()

	// WHEN getting payment fails
	response, err := suite.controller.GetPaymentByOrderId(orderId)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	assert.Nil(suite.T(), response)
	suite.mockGetPaymentUseCase.AssertExpectations(suite.T())
}

func (suite *PaymentControllerTestSuite) Test_UpdatePaymentStatus_WithValidData_ShouldSucceed() {
	// GIVEN valid order ID and status
	orderId := uint(1)
	status := "Approved"

	suite.mockUpdatePaymentUseCase.EXPECT().
		Execute(mock.Anything).
		Return(nil).
		Once()

	// WHEN updating payment status
	err := suite.controller.UpdatePaymentStatus(orderId, status)

	// THEN operation should succeed
	assert.NoError(suite.T(), err)
	suite.mockUpdatePaymentUseCase.AssertExpectations(suite.T())
}

func (suite *PaymentControllerTestSuite) Test_UpdatePaymentStatus_WithError_ShouldReturnError() {
	// GIVEN an order ID and status
	orderId := uint(1)
	status := "Approved"
	expectedError := errors.New("update failed")

	suite.mockUpdatePaymentUseCase.EXPECT().
		Execute(mock.Anything).
		Return(expectedError).
		Once()

	// WHEN update fails
	err := suite.controller.UpdatePaymentStatus(orderId, status)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockUpdatePaymentUseCase.AssertExpectations(suite.T())
}
