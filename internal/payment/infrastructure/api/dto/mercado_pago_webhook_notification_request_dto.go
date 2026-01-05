package dto

type MercadoPagoWebhookNotificationRequestDTO struct {
	Id       string `json:"id"`
	Topic    string `json:"topic"`
	Resource string `json:"resource"`
}
