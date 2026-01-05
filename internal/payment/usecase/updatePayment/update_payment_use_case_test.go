package updatepayment_test

import (
	"errors"
	"testing"

	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
	updatepayment "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/updatePayment"
	mockRepositories "github.com/abattassini/tc-fiap-payment/mocks/payment/domain/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UpdatePaymentUseCaseTestSuite struct {
	suite.Suite
	mockRepository *mockRepositories.MockPaymentRepository
	useCase        updatepayment.UpdatePaymentUseCase
}

func (suite *UpdatePaymentUseCaseTestSuite) SetupTest() {
	suite.mockRepository = mockRepositories.NewMockPaymentRepository(suite.T())
	suite.useCase = updatepayment.NewUpdatePaymentUseCaseImpl(suite.mockRepository)
}

func TestUpdatePaymentUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UpdatePaymentUseCaseTestSuite))
}

func (suite *UpdatePaymentUseCaseTestSuite) Test_UpdatePaymentStatus_WithExistingPayment_ShouldUpdate() {
	// GIVEN an existing payment
	orderId := uint(1)
	command := commands.NewUpdatePaymentStatusCommand(orderId, "Approved")

	payment := &entities.Payment{
		ID:      1,
		OrderId: orderId,
		Status:  "pending",
	}

	suite.mockRepository.EXPECT().
		GetPaymentByOrderId(orderId).
		Return(payment, nil).
		Once()

	suite.mockRepository.EXPECT().
		UpdatePayment(payment).
		Return(nil).
		Once()

	// WHEN updating the payment status
	err := suite.useCase.Execute(command)

	// THEN the operation should complete without errors
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Approved", payment.Status)
	suite.mockRepository.AssertExpectations(suite.T())
}

func (suite *UpdatePaymentUseCaseTestSuite) Test_UpdatePaymentStatus_WithNonExistentPayment_ShouldReturnError() {
	// GIVEN a non-existent payment
	orderId := uint(999)
	command := commands.NewUpdatePaymentStatusCommand(orderId, "Approved")

	expectedError := errors.New("payment not found")

	suite.mockRepository.EXPECT().
		GetPaymentByOrderId(orderId).
		Return(nil, expectedError).
		Once()

	// WHEN updating the payment status
	err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockRepository.AssertExpectations(suite.T())
}

func (suite *UpdatePaymentUseCaseTestSuite) Test_UpdatePaymentStatus_WithRepositoryFailure_ShouldReturnError() {
	// GIVEN an existing payment
	orderId := uint(1)
	command := commands.NewUpdatePaymentStatusCommand(orderId, "Declined")

	payment := &entities.Payment{
		ID:      1,
		OrderId: orderId,
		Status:  "pending",
	}

	expectedError := errors.New("database error")

	suite.mockRepository.EXPECT().
		GetPaymentByOrderId(orderId).
		Return(payment, nil).
		Once()

	suite.mockRepository.EXPECT().
		UpdatePayment(payment).
		Return(expectedError).
		Once()

	// WHEN updating the payment status
	err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockRepository.AssertExpectations(suite.T())
}
