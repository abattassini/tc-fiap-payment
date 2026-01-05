package addpayment

import (
	"context"
	"fmt"
	"os"

	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/repositories"
	"github.com/abattassini/tc-fiap-payment/internal/payment/gateways"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/clients"
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
)

var (
	_ AddPaymentUseCase = (*AddPaymentUseCaseImpl)(nil)
)

type AddPaymentUseCaseImpl struct {
	mercadoPagoGateway gateways.MercadoPagoGateway
	orderClient        clients.OrderClient
	paymentRepository  repositories.PaymentRepository
}

func NewAddPaymentUseCaseImpl(
	mercadoPagoGateway gateways.MercadoPagoGateway,
	orderClient clients.OrderClient,
	paymentRepository repositories.PaymentRepository) *AddPaymentUseCaseImpl {
	return &AddPaymentUseCaseImpl{
		mercadoPagoGateway: mercadoPagoGateway,
		orderClient:        orderClient,
		paymentRepository:  paymentRepository,
	}
}

func (u *AddPaymentUseCaseImpl) Execute(command *commands.AddPaymentCommand) (string, error) {
	paymentEntity := entities.Payment{
		OrderId: command.OrderId,
		Total:   command.Total,
		Type:    command.Type,
		Status:  "pending",
	}

	paymentResult, err := u.paymentRepository.AddPayment(&paymentEntity)
	if err != nil {
		return "", err
	}

	order, err := u.orderClient.GetOrder(paymentResult.OrderId)
	if err != nil {
		return "", err
	}

	var items []dto.Item

	for _, product := range order.Products {
		items = append(items, dto.Item{
			SKUNumber:   fmt.Sprint(product.ProductId),
			Category:    fmt.Sprint(product.Category),
			Title:       product.Name,
			Description: product.Description,
			UnitPrice:   product.Price,
			Quantity:    int(product.Quantity),
			TotalAmount: float32(product.Quantity) * product.Price,
		})
	}

	qrCodeResponse, err := u.mercadoPagoGateway.GenerateQRCode(context.Background(), dto.CreateQRCodeDTO{
		ExternalReference: fmt.Sprintf("order-%d", order.ID),
		Title:             "Fiap",
		Description:       "Fiap",
		NotificationURL:   os.Getenv("MERCADO_PAGO_WEBHOOK_CALLBACK_URL"),
		TotalAmount:       order.TotalAmount,
		Items:             items,
	})
	if err != nil {
		return "", err
	}

	return qrCodeResponse.QRData, nil
}
