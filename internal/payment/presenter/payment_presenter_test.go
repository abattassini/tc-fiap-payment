package presenter_test

import (
	"testing"
	"time"

	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
	"github.com/abattassini/tc-fiap-payment/internal/payment/presenter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PaymentPresenterTestSuite struct {
	suite.Suite
	presenter presenter.PaymentPresenter
}

func (suite *PaymentPresenterTestSuite) SetupTest() {
	suite.presenter = presenter.NewPaymentPresenterImpl()
}

func TestPaymentPresenterTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentPresenterTestSuite))
}

func (suite *PaymentPresenterTestSuite) Test_Present_WithValidPayment_ShouldReturnDTO() {
	// GIVEN a payment entity
	payment := &entities.Payment{
		ID:        1,
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		OrderId:   123,
		Total:     250.75,
		Type:      "credit_card",
		Status:    "Approved",
	}

	// WHEN presenting the payment
	dto := suite.presenter.Present(payment)

	// THEN the DTO should contain all payment data
	assert.NotNil(suite.T(), dto)
	assert.Equal(suite.T(), uint(1), dto.ID)
	assert.Equal(suite.T(), uint(123), dto.OrderId)
	assert.Equal(suite.T(), float32(250.75), dto.Total)
	assert.Equal(suite.T(), "credit_card", dto.Type)
	assert.Equal(suite.T(), "Approved", dto.Status)
	assert.Equal(suite.T(), payment.CreatedAt, dto.CreatedAt)
}
