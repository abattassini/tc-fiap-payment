package controller

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
	handleWebhookUseCase "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/handleWebhook"
)

var (
	_ PaymentWebhookController = (*PaymentWebhookControllerImpl)(nil)
)

type PaymentWebhookControllerImpl struct {
	handleWebhookUseCase handleWebhookUseCase.HandleWebhookUseCase
}

func NewPaymentWebhookControllerImpl(handleWebhookUseCase handleWebhookUseCase.HandleWebhookUseCase) *PaymentWebhookControllerImpl {
	return &PaymentWebhookControllerImpl{handleWebhookUseCase: handleWebhookUseCase}
}

func (c *PaymentWebhookControllerImpl) HandleWebhook(mercadoPagoWebhookRequest *dto.MercadoPagoWebhookNotificationRequestDTO) error {
	status := GetStatusFromString(mercadoPagoWebhookRequest.Topic)

	command := commands.HandleWebhookCommand{
		Id:     mercadoPagoWebhookRequest.Id,
		Status: status,
	}
	return c.handleWebhookUseCase.Execute(command)
}

func GetStatusFromString(status string) string {
	switch status {
	case "payment.created", "payment.updated":
		return "Approved"
	default:
		return "Declined"
	}
}
