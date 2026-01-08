package handlewebhook_test

import (
	"errors"
	"testing"

	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
	handlewebhook "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/handleWebhook"
	mockClients "github.com/abattassini/tc-fiap-payment/mocks/payment/infrastructure/clients"
	mockUpdatePayment "github.com/abattassini/tc-fiap-payment/mocks/payment/usecase/updatePayment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HandleWebhookUseCaseTestSuite struct {
	suite.Suite
	mockUpdatePaymentUseCase *mockUpdatePayment.MockUpdatePaymentUseCase
	mockOrderClient          *mockClients.MockOrderClient
	useCase                  handlewebhook.HandleWebhookUseCase
}

func (suite *HandleWebhookUseCaseTestSuite) SetupTest() {
	suite.mockUpdatePaymentUseCase = mockUpdatePayment.NewMockUpdatePaymentUseCase(suite.T())
	suite.mockOrderClient = mockClients.NewMockOrderClient(suite.T())
	suite.useCase = handlewebhook.NewHandleWebhookUseCaseImpl(
		suite.mockUpdatePaymentUseCase,
		suite.mockOrderClient,
	)
}

func TestHandleWebhookUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(HandleWebhookUseCaseTestSuite))
}

func (suite *HandleWebhookUseCaseTestSuite) Test_HandleWebhook_WithApprovedPayment_ShouldUpdateOrderStatus() {
	// GIVEN an approved payment webhook
	command := commands.HandleWebhookCommand{
		Id:     "1",
		Status: "Approved",
	}

	suite.mockUpdatePaymentUseCase.EXPECT().
		Execute(mock.MatchedBy(func(cmd *commands.UpdatePaymentStatusCommand) bool {
			return cmd.OrderId == 1 && cmd.Status == "Approved"
		})).
		Return(nil).
		Once()

	suite.mockOrderClient.EXPECT().
		UpdateOrderStatus(uint(1), int(2)).
		Return(nil).
		Once()

	// WHEN handling webhook
	err := suite.useCase.Execute(command)

	// THEN payment and order should be updated
	assert.NoError(suite.T(), err)
	suite.mockUpdatePaymentUseCase.AssertExpectations(suite.T())
	suite.mockOrderClient.AssertExpectations(suite.T())
}

func (suite *HandleWebhookUseCaseTestSuite) Test_HandleWebhook_WithDeclinedPayment_ShouldOnlyUpdatePayment() {
	// GIVEN a declined payment webhook
	command := commands.HandleWebhookCommand{
		Id:     "1",
		Status: "Declined",
	}

	suite.mockUpdatePaymentUseCase.EXPECT().
		Execute(mock.MatchedBy(func(cmd *commands.UpdatePaymentStatusCommand) bool {
			return cmd.OrderId == 1 && cmd.Status == "Declined"
		})).
		Return(nil).
		Once()

	// WHEN handling webhook with non-approved status
	err := suite.useCase.Execute(command)

	// THEN only payment should be updated
	assert.NoError(suite.T(), err)
	suite.mockUpdatePaymentUseCase.AssertExpectations(suite.T())
}

func (suite *HandleWebhookUseCaseTestSuite) Test_HandleWebhook_WithInvalidOrderId_ShouldReturnError() {
	// GIVEN a webhook with invalid order ID
	command := commands.HandleWebhookCommand{
		Id:     "invalid",
		Status: "Approved",
	}

	// WHEN handling webhook with invalid ID
	err := suite.useCase.Execute(command)

	// THEN error should be returned
	assert.Error(suite.T(), err)
}

func (suite *HandleWebhookUseCaseTestSuite) Test_HandleWebhook_WithUpdatePaymentError_ShouldReturnError() {
	// GIVEN an approved payment webhook
	command := commands.HandleWebhookCommand{
		Id:     "1",
		Status: "Approved",
	}

	expectedError := errors.New("payment update failed")

	suite.mockUpdatePaymentUseCase.EXPECT().
		Execute(mock.Anything).
		Return(expectedError).
		Once()

	// WHEN payment update fails
	err := suite.useCase.Execute(command)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockUpdatePaymentUseCase.AssertExpectations(suite.T())
}

func (suite *HandleWebhookUseCaseTestSuite) Test_HandleWebhook_WithOrderClientError_ShouldReturnError() {
	// GIVEN an approved payment webhook
	command := commands.HandleWebhookCommand{
		Id:     "1",
		Status: "Approved",
	}

	expectedError := errors.New("order service unavailable")

	suite.mockUpdatePaymentUseCase.EXPECT().
		Execute(mock.Anything).
		Return(nil).
		Once()

	suite.mockOrderClient.EXPECT().
		UpdateOrderStatus(uint(1), int(2)).
		Return(expectedError).
		Once()

	// WHEN order update fails
	err := suite.useCase.Execute(command)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to update order status in Order Service")
	assert.Contains(suite.T(), err.Error(), expectedError.Error())
	suite.mockUpdatePaymentUseCase.AssertExpectations(suite.T())
	suite.mockOrderClient.AssertExpectations(suite.T())
}
