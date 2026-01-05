package getpaymentstatus_test

import (
	"errors"
	"testing"

	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
	getpaymentstatus "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/getPaymentStatus"
	mockRepositories "github.com/abattassini/tc-fiap-payment/mocks/payment/domain/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetPaymentStatusUseCaseTestSuite struct {
	suite.Suite
	mockRepository *mockRepositories.MockPaymentRepository
	useCase        getpaymentstatus.GetPaymentStatusUseCase
}

func (suite *GetPaymentStatusUseCaseTestSuite) SetupTest() {
	suite.mockRepository = mockRepositories.NewMockPaymentRepository(suite.T())
	suite.useCase = getpaymentstatus.NewGetPaymentStatusUseCaseImpl(suite.mockRepository)
}

func TestGetPaymentStatusUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(GetPaymentStatusUseCaseTestSuite))
}

func (suite *GetPaymentStatusUseCaseTestSuite) Test_GetPaymentStatus_WithExistingPayment_ShouldReturnStatus() {
	// GIVEN an order with a payment
	orderId := uint(123)
	command := commands.NewGetPaymentStatusCommand(orderId)

	payment := &entities.Payment{
		ID:      1,
		OrderId: orderId,
		Status:  "pending",
	}

	suite.mockRepository.EXPECT().
		GetPaymentByOrderId(orderId).
		Return(payment, nil).
		Once()

	// WHEN the payment status is requested
	status, err := suite.useCase.Execute(command)

	// THEN the status should be returned successfully
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "pending", status)
	suite.mockRepository.AssertExpectations(suite.T())
}

func (suite *GetPaymentStatusUseCaseTestSuite) Test_GetPaymentStatus_WithNonExistentPayment_ShouldReturnError() {
	// GIVEN an order without a payment
	orderId := uint(999)
	command := commands.NewGetPaymentStatusCommand(orderId)

	expectedError := errors.New("payment not found")

	suite.mockRepository.EXPECT().
		GetPaymentByOrderId(orderId).
		Return(nil, expectedError).
		Once()

	// WHEN the payment status is requested
	status, err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "", status)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockRepository.AssertExpectations(suite.T())
}
