package updatepayment

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
)

type UpdatePaymentUseCase interface {
	Execute(command *commands.UpdatePaymentStatusCommand) error
}
