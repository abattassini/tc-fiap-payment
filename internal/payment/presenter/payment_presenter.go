package presenter

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
)

type PaymentPresenter interface {
	Present(payment *entities.Payment) *dto.GetPaymentResponseDto
}
