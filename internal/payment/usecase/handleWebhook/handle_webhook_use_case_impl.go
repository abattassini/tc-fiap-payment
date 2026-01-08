package handlewebhook

import (
	"fmt"
	"strconv"

	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/clients"
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
	updatePaymentUseCase "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/updatePayment"
)

var (
	_ HandleWebhookUseCase = (*HandleWebhookUseCaseImpl)(nil)
)

type HandleWebhookUseCaseImpl struct {
	updatePaymentUseCase updatePaymentUseCase.UpdatePaymentUseCase
	orderClient          clients.OrderClient
}

func NewHandleWebhookUseCaseImpl(
	updatePaymentUseCase updatePaymentUseCase.UpdatePaymentUseCase,
	orderClient clients.OrderClient) *HandleWebhookUseCaseImpl {
	return &HandleWebhookUseCaseImpl{
		updatePaymentUseCase: updatePaymentUseCase,
		orderClient:          orderClient,
	}
}

func (u *HandleWebhookUseCaseImpl) Execute(command commands.HandleWebhookCommand) error {
	orderId, err := strconv.ParseUint(command.Id, 10, 32)
	if err != nil {
		return err
	}

	updatePayment := commands.UpdatePaymentStatusCommand{
		OrderId: uint(orderId),
		Status:  command.Status,
	}

	err = u.updatePaymentUseCase.Execute(&updatePayment)
	if err != nil {
		return err
	}

	if command.Status == "Approved" {
		// Update order status to "Preparing" (status=2) when payment is approved
		err = u.orderClient.UpdateOrderStatus(uint(orderId), 2)
		if err != nil {
			println("ERROR: Failed to update order status in Order Service:", err.Error())
			return fmt.Errorf("failed to update order status in Order Service: %w", err)
		}
		println("Successfully updated order", orderId, "status to 'Preparing' (status=2)")
	}

	return nil
}
