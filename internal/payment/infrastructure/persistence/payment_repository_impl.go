package persistence

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/repositories"
	"gorm.io/gorm"
)

var (
	_ repositories.PaymentRepository = (*PaymentRepositoryImpl)(nil)
)

type PaymentRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRepositoryImpl(db *gorm.DB) *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{db: db}
}

func (r *PaymentRepositoryImpl) AddPayment(payment *entities.Payment) (*entities.Payment, error) {
	if err := r.db.Create(payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *PaymentRepositoryImpl) GetPaymentByOrderId(orderId uint) (*entities.Payment, error) {
	payment := &entities.Payment{}
	if err := r.db.
		Where("order_id = ?", orderId).
		First(payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *PaymentRepositoryImpl) UpdatePayment(payment *entities.Payment) error {
	return r.db.Save(payment).Error
}
