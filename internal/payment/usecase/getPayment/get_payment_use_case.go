package getpayment

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
)

type GetPaymentUseCase interface {
	Execute(command *commands.GetPaymentCommand) (*entities.Payment, error)
}
