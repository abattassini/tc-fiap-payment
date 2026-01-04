package commands

type AddPaymentCommand struct {
	OrderId uint
	Total   float32
	Type    string
}

func NewAddPaymentCommand(orderId uint, total float32, type_ string) *AddPaymentCommand {
	return &AddPaymentCommand{
		OrderId: orderId,
		Total:   total,
		Type:    type_,
	}
}
