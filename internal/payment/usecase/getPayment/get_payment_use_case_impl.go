package getpayment

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/repositories"
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
)

var (
	_ GetPaymentUseCase = (*GetPaymentUseCaseImpl)(nil)
)

type GetPaymentUseCaseImpl struct {
	paymentRepository repositories.PaymentRepository
}

func NewGetPaymentUseCaseImpl(paymentRepository repositories.PaymentRepository) *GetPaymentUseCaseImpl {
	return &GetPaymentUseCaseImpl{paymentRepository: paymentRepository}
}

func (u *GetPaymentUseCaseImpl) Execute(command *commands.GetPaymentCommand) (*entities.Payment, error) {
	payment, err := u.paymentRepository.GetPaymentByOrderId(command.OrderId)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
