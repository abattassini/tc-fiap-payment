package dto

type OrderResponseDto struct {
	ID          uint                `json:"id"`
	TotalAmount float32             `json:"total_amount"`
	Products    []*OrderProductDto  `json:"products"`
}

type OrderProductDto struct {
	ProductId   uint    `json:"product_id"`
	Price       float32 `json:"price"`
	Quantity    uint    `json:"quantity"`
	Name        string  `json:"name"`
	ImageLink   string  `json:"image_link"`
	Description string  `json:"description"`
	Category    int     `json:"category"`
}
