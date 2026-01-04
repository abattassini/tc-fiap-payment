package dto

type QRCodeResponseDto struct {
	QRData         string `json:"qr_data"`
	InStoreOrderId string `json:"in_store_order_id"`
}
