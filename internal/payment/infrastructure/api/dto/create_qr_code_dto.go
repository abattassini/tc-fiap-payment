package dto

type CreateQRCodeDTO struct {
	ExternalReference string  `json:"external_reference"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	NotificationURL   string  `json:"notification_url"`
	TotalAmount       float32 `json:"total_amount"`
	Items             []Item  `json:"items"`
}

type Item struct {
	SKUNumber   string  `json:"sku_number"`
	Category    string  `json:"category"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
	UnitMeasure string  `json:"unit_measure"`
	TotalAmount float32 `json:"total_amount"`
}
