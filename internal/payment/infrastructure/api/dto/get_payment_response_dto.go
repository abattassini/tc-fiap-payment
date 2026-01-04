package dto

import "time"

type GetPaymentResponseDto struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	OrderId   uint      `json:"order_id"`
	Total     float32   `json:"total"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
}
