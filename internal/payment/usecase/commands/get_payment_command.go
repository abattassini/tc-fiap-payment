package commands

type GetPaymentCommand struct {
	OrderId uint
}

func NewGetPaymentCommand(orderId uint) *GetPaymentCommand {
	return &GetPaymentCommand{
		OrderId: orderId,
	}
}
