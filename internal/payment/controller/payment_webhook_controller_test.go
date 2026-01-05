package controller_test

import (
	"errors"
	"testing"

	"github.com/abattassini/tc-fiap-payment/internal/payment/controller"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	mockHandleWebhook "github.com/abattassini/tc-fiap-payment/mocks/payment/usecase/handleWebhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PaymentWebhookControllerTestSuite struct {
	suite.Suite
	mockHandleWebhookUseCase *mockHandleWebhook.MockHandleWebhookUseCase
	controller               controller.PaymentWebhookController
}

func (suite *PaymentWebhookControllerTestSuite) SetupTest() {
	suite.mockHandleWebhookUseCase = mockHandleWebhook.NewMockHandleWebhookUseCase(suite.T())
	suite.controller = controller.NewPaymentWebhookControllerImpl(suite.mockHandleWebhookUseCase)
}

func TestPaymentWebhookControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentWebhookControllerTestSuite))
}

func (suite *PaymentWebhookControllerTestSuite) Test_HandleWebhook_WithPaymentCreated_ShouldProcessAsApproved() {
	// GIVEN a payment.created webhook
	request := &dto.MercadoPagoWebhookNotificationRequestDTO{
		Id:    "1",
		Topic: "payment.created",
	}

	suite.mockHandleWebhookUseCase.EXPECT().
		Execute(mock.MatchedBy(func(cmd interface{}) bool {
			return true
		})).
		Return(nil).
		Once()

	// WHEN handling webhook
	err := suite.controller.HandleWebhook(request)

	// THEN operation should succeed
	assert.NoError(suite.T(), err)
	suite.mockHandleWebhookUseCase.AssertExpectations(suite.T())
}

func (suite *PaymentWebhookControllerTestSuite) Test_HandleWebhook_WithPaymentUpdated_ShouldProcessAsApproved() {
	// GIVEN a payment.updated webhook
	request := &dto.MercadoPagoWebhookNotificationRequestDTO{
		Id:    "1",
		Topic: "payment.updated",
	}

	suite.mockHandleWebhookUseCase.EXPECT().
		Execute(mock.Anything).
		Return(nil).
		Once()

	// WHEN handling webhook
	err := suite.controller.HandleWebhook(request)

	// THEN operation should succeed
	assert.NoError(suite.T(), err)
	suite.mockHandleWebhookUseCase.AssertExpectations(suite.T())
}

func (suite *PaymentWebhookControllerTestSuite) Test_HandleWebhook_WithOtherTopic_ShouldProcessAsDeclined() {
	// GIVEN a webhook with other topic
	request := &dto.MercadoPagoWebhookNotificationRequestDTO{
		Id:    "1",
		Topic: "payment.failed",
	}

	suite.mockHandleWebhookUseCase.EXPECT().
		Execute(mock.Anything).
		Return(nil).
		Once()

	// WHEN handling webhook
	err := suite.controller.HandleWebhook(request)

	// THEN operation should succeed
	assert.NoError(suite.T(), err)
	suite.mockHandleWebhookUseCase.AssertExpectations(suite.T())
}

func (suite *PaymentWebhookControllerTestSuite) Test_HandleWebhook_WithError_ShouldReturnError() {
	// GIVEN a webhook request
	request := &dto.MercadoPagoWebhookNotificationRequestDTO{
		Id:    "1",
		Topic: "payment.created",
	}

	expectedError := errors.New("webhook processing failed")

	suite.mockHandleWebhookUseCase.EXPECT().
		Execute(mock.Anything).
		Return(expectedError).
		Once()

	// WHEN webhook processing fails
	err := suite.controller.HandleWebhook(request)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockHandleWebhookUseCase.AssertExpectations(suite.T())
}

func (suite *PaymentWebhookControllerTestSuite) Test_GetStatusFromString_WithPaymentCreated_ShouldReturnApproved() {
	// WHEN getting status from payment.created
	status := controller.GetStatusFromString("payment.created")

	// THEN should return Approved
	assert.Equal(suite.T(), "Approved", status)
}

func (suite *PaymentWebhookControllerTestSuite) Test_GetStatusFromString_WithPaymentUpdated_ShouldReturnApproved() {
	// WHEN getting status from payment.updated
	status := controller.GetStatusFromString("payment.updated")

	// THEN should return Approved
	assert.Equal(suite.T(), "Approved", status)
}

func (suite *PaymentWebhookControllerTestSuite) Test_GetStatusFromString_WithOtherTopic_ShouldReturnDeclined() {
	// WHEN getting status from other topic
	status := controller.GetStatusFromString("payment.failed")

	// THEN should return Declined
	assert.Equal(suite.T(), "Declined", status)
}
