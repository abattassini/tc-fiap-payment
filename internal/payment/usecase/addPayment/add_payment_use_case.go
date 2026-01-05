package addpayment

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
)

type AddPaymentUseCase interface {
	Execute(command *commands.AddPaymentCommand) (string, error)
}
