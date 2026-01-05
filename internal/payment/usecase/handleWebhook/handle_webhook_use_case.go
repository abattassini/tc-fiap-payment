package handlewebhook

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
)

type HandleWebhookUseCase interface {
	Execute(command commands.HandleWebhookCommand) error
}
