package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	paymentController "github.com/abattassini/tc-fiap-payment/internal/payment/controller"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	"github.com/go-chi/chi/v5"
)

type PaymentApiController struct {
	paymentController paymentController.PaymentController
}

func NewPaymentApiController(paymentService paymentController.PaymentController) *PaymentApiController {
	return &PaymentApiController{paymentController: paymentService}
}

func (c *PaymentApiController) RegisterRoutes(r chi.Router) {
	prefix := "/v1/payment"
	r.Post(prefix, c.CreatePayment)
	r.Get(prefix+"/{orderId}/status", c.GetPaymentStatusByOrderId)
	r.Get(prefix+"/{orderId}", c.GetPaymentByOrderId)
}

func (c *PaymentApiController) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var request dto.AddPaymentRequestDto

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	paymentCode, err := c.paymentController.CreatePayment(&request)
	if err != nil {
		// Log the actual error for debugging
		println("Error creating payment:", err.Error())
		http.Error(w, fmt.Sprintf("Error processing request: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(paymentCode)
}

func (c *PaymentApiController) GetPaymentStatusByOrderId(w http.ResponseWriter, r *http.Request) {
	orderId, err := getOrderIDFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := c.paymentController.GetPaymentStatusByOrderId(orderId)
	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

func (c *PaymentApiController) GetPaymentByOrderId(w http.ResponseWriter, r *http.Request) {
	orderId, err := getOrderIDFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment, err := c.paymentController.GetPaymentByOrderId(orderId)
	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

func getOrderIDFromPath(r *http.Request) (uint, error) {
	vars := chi.URLParam(r, "orderId")
	id, err := strconv.ParseUint(vars, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
