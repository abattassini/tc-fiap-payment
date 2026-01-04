package getpaymentstatus

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/repositories"
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
)

var (
	_ GetPaymentStatusUseCase = (*GetPaymentStatusUseCaseImpl)(nil)
)

type GetPaymentStatusUseCaseImpl struct {
	paymentRepository repositories.PaymentRepository
}

func NewGetPaymentStatusUseCaseImpl(paymentRepository repositories.PaymentRepository) *GetPaymentStatusUseCaseImpl {
	return &GetPaymentStatusUseCaseImpl{paymentRepository: paymentRepository}
}

func (u *GetPaymentStatusUseCaseImpl) Execute(command *commands.GetPaymentStatusCommand) (string, error) {
	payment, err := u.paymentRepository.GetPaymentByOrderId(command.OrderId)
	if err != nil {
		return "", err
	}

	return payment.Status, nil
}
