package controller

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
)

type PaymentController interface {
	CreatePayment(addPaymentRequest *dto.AddPaymentRequestDto) (string, error)
	GetPaymentStatusByOrderId(orderId uint) (string, error)
	GetPaymentByOrderId(orderId uint) (*dto.GetPaymentResponseDto, error)
	UpdatePaymentStatus(orderId uint, status string) error
}
