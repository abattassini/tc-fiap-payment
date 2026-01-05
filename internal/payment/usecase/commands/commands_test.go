package commands_test

import (
	"testing"

	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
	"github.com/stretchr/testify/assert"
)

func TestNewAddPaymentCommand(t *testing.T) {
	// GIVEN payment data
	orderId := uint(1)
	total := float32(100.50)
	paymentType := "QRCode"

	// WHEN creating command
	cmd := commands.NewAddPaymentCommand(orderId, total, paymentType)

	// THEN command should be created correctly
	assert.NotNil(t, cmd)
	assert.Equal(t, orderId, cmd.OrderId)
	assert.Equal(t, total, cmd.Total)
	assert.Equal(t, paymentType, cmd.Type)
}

func TestNewGetPaymentCommand(t *testing.T) {
	// GIVEN an order ID
	orderId := uint(1)

	// WHEN creating command
	cmd := commands.NewGetPaymentCommand(orderId)

	// THEN command should be created correctly
	assert.NotNil(t, cmd)
	assert.Equal(t, orderId, cmd.OrderId)
}

func TestNewGetPaymentStatusCommand(t *testing.T) {
	// GIVEN an order ID
	orderId := uint(1)

	// WHEN creating command
	cmd := commands.NewGetPaymentStatusCommand(orderId)

	// THEN command should be created correctly
	assert.NotNil(t, cmd)
	assert.Equal(t, orderId, cmd.OrderId)
}

func TestNewUpdatePaymentStatusCommand(t *testing.T) {
	// GIVEN order ID and status
	orderId := uint(1)
	status := "Approved"

	// WHEN creating command
	cmd := commands.NewUpdatePaymentStatusCommand(orderId, status)

	// THEN command should be created correctly
	assert.NotNil(t, cmd)
	assert.Equal(t, orderId, cmd.OrderId)
	assert.Equal(t, status, cmd.Status)
}

func TestHandleWebhookCommand(t *testing.T) {
	// GIVEN webhook data
	id := "123"
	status := "Approved"

	// WHEN creating command
	cmd := commands.HandleWebhookCommand{
		Id:     id,
		Status: status,
	}

	// THEN command should have correct values
	assert.Equal(t, id, cmd.Id)
	assert.Equal(t, status, cmd.Status)
}
