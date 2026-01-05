package presenter

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
)

var (
	_ PaymentPresenter = (*PaymentPresenterImpl)(nil)
)

type PaymentPresenterImpl struct{}

func NewPaymentPresenterImpl() *PaymentPresenterImpl {
	return &PaymentPresenterImpl{}
}

func (p *PaymentPresenterImpl) Present(payment *entities.Payment) *dto.GetPaymentResponseDto {
	return &dto.GetPaymentResponseDto{
		ID:        payment.ID,
		CreatedAt: payment.CreatedAt,
		OrderId:   payment.OrderId,
		Total:     payment.Total,
		Type:      payment.Type,
		Status:    payment.Status,
	}
}
