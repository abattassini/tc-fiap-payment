package getpaymentstatus

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
)

type GetPaymentStatusUseCase interface {
	Execute(command *commands.GetPaymentStatusCommand) (string, error)
}
