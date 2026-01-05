package handlewebhook

import (
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
		// TODO: Replace mock with real Order Service call once the service is running
		// Currently mocking the Order Service UpdateOrderStatus to avoid dependency on external service
		// Original code:
		// err = u.orderClient.UpdateOrderStatus(uint(orderId), 2)
		// if err != nil {
		//     return err
		// }

		// Mock: Simulating successful order status update
		// In production, this would call the Order Service to update status to "Preparing" (status=2)
		_ = uint(orderId) // Just to avoid unused variable warning
	}

	return nil
}
