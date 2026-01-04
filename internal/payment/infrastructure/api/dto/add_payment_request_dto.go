package dto

type AddPaymentRequestDto struct {
	OrderId uint    `json:"orderId"`
	Total   float32 `json:"total"`
	Type    string  `json:"type"`
}
