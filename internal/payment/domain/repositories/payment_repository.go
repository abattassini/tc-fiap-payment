package repositories

import "github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"

type PaymentRepository interface {
	AddPayment(payment *entities.Payment) (*entities.Payment, error)
	GetPaymentByOrderId(orderId uint) (*entities.Payment, error)
	UpdatePayment(payment *entities.Payment) error
}
