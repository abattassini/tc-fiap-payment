package commands

type GetPaymentStatusCommand struct {
	OrderId uint
}

func NewGetPaymentStatusCommand(orderId uint) *GetPaymentStatusCommand {
	return &GetPaymentStatusCommand{
		OrderId: orderId,
	}
}
