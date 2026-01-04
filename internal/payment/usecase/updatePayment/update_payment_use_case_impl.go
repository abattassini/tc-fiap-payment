package updatepayment

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/repositories"
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
)

var (
	_ UpdatePaymentUseCase = (*UpdatePaymentUseCaseImpl)(nil)
)

type UpdatePaymentUseCaseImpl struct {
	paymentRepository repositories.PaymentRepository
}

func NewUpdatePaymentUseCaseImpl(paymentRepository repositories.PaymentRepository) *UpdatePaymentUseCaseImpl {
	return &UpdatePaymentUseCaseImpl{paymentRepository: paymentRepository}
}

func (u *UpdatePaymentUseCaseImpl) Execute(command *commands.UpdatePaymentStatusCommand) error {
	payment, err := u.paymentRepository.GetPaymentByOrderId(command.OrderId)
	if err != nil {
		return err
	}

	payment.Status = command.Status
	return u.paymentRepository.UpdatePayment(payment)
}
