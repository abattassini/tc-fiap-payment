package controller

import "github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"

type PaymentWebhookController interface {
	HandleWebhook(mercadoPagoWebhookRequest *dto.MercadoPagoWebhookNotificationRequestDTO) error
}
