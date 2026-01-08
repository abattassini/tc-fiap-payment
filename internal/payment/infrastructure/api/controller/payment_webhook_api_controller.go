package controller

import (
	"encoding/json"
	"net/http"

	"github.com/abattassini/tc-fiap-payment/internal/payment/controller"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	"github.com/go-chi/chi/v5"
)

type PaymentWebhookApiController struct {
	paymentWebhookController controller.PaymentWebhookController
}

func NewPaymentWebhookApiController(paymentWebhookController controller.PaymentWebhookController) *PaymentWebhookApiController {
	return &PaymentWebhookApiController{paymentWebhookController: paymentWebhookController}
}

func (c *PaymentWebhookApiController) RegisterRoutes(r chi.Router) {
	prefix := "/payment/webhooks"
	r.Post(prefix+"/notify", c.HandlePaymentNotification)
}

func (c *PaymentWebhookApiController) HandlePaymentNotification(w http.ResponseWriter, r *http.Request) {
	var request dto.MercadoPagoWebhookNotificationRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := c.paymentWebhookController.HandleWebhook(&request)
	if err != nil {
		http.Error(w, "Error processing webhook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
