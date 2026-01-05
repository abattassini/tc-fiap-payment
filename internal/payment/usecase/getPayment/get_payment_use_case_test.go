package getpayment_test

import (
	"errors"
	"testing"
	"time"

	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
	getpayment "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/getPayment"
	mockRepositories "github.com/abattassini/tc-fiap-payment/mocks/payment/domain/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetPaymentUseCaseTestSuite struct {
	suite.Suite
	mockRepository *mockRepositories.MockPaymentRepository
	useCase        getpayment.GetPaymentUseCase
}

func (suite *GetPaymentUseCaseTestSuite) SetupTest() {
	suite.mockRepository = mockRepositories.NewMockPaymentRepository(suite.T())
	suite.useCase = getpayment.NewGetPaymentUseCaseImpl(suite.mockRepository)
}

func TestGetPaymentUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(GetPaymentUseCaseTestSuite))
}

func (suite *GetPaymentUseCaseTestSuite) Test_GetPayment_WithExistingPayment_ShouldReturnPayment() {
	// GIVEN an existing payment
	orderId := uint(1)
	command := commands.NewGetPaymentCommand(orderId)

	expectedPayment := &entities.Payment{
		ID:        1,
		CreatedAt: time.Now(),
		OrderId:   orderId,
		Total:     100.50,
		Type:      "credit_card",
		Status:    "pending",
	}

	suite.mockRepository.EXPECT().
		GetPaymentByOrderId(orderId).
		Return(expectedPayment, nil).
		Once()

	// WHEN getting the payment
	payment, err := suite.useCase.Execute(command)

	// THEN the payment should be returned without errors
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), payment)
	assert.Equal(suite.T(), expectedPayment.ID, payment.ID)
	assert.Equal(suite.T(), expectedPayment.OrderId, payment.OrderId)
	assert.Equal(suite.T(), expectedPayment.Total, payment.Total)
	assert.Equal(suite.T(), expectedPayment.Status, payment.Status)
	suite.mockRepository.AssertExpectations(suite.T())
}

func (suite *GetPaymentUseCaseTestSuite) Test_GetPayment_WithNonExistentPayment_ShouldReturnError() {
	// GIVEN a non-existent payment
	orderId := uint(999)
	command := commands.NewGetPaymentCommand(orderId)

	expectedError := errors.New("payment not found")

	suite.mockRepository.EXPECT().
		GetPaymentByOrderId(orderId).
		Return(nil, expectedError).
		Once()

	// WHEN getting the payment
	payment, err := suite.useCase.Execute(command)

	// THEN an error should be returned
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), payment)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockRepository.AssertExpectations(suite.T())
}
