package commands

type UpdatePaymentStatusCommand struct {
	OrderId uint
	Status  string
}

func NewUpdatePaymentStatusCommand(orderId uint, status string) *UpdatePaymentStatusCommand {
	return &UpdatePaymentStatusCommand{
		OrderId: orderId,
		Status:  status,
	}
}
