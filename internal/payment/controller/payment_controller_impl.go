package controller

import (
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	paymentPresenter "github.com/abattassini/tc-fiap-payment/internal/payment/presenter"
	addPayment "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/addPayment"
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
	getpayment "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/getPayment"
	getpaymentstatus "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/getPaymentStatus"
	updatepayment "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/updatePayment"
)

var (
	_ PaymentController = (*PaymentControllerImpl)(nil)
)

type PaymentControllerImpl struct {
	presenter               paymentPresenter.PaymentPresenter
	getPaymentUseCase       getpayment.GetPaymentUseCase
	getPaymentStatusUseCase getpaymentstatus.GetPaymentStatusUseCase
	updatePaymentUseCase    updatepayment.UpdatePaymentUseCase
	addPaymentUseCase       addPayment.AddPaymentUseCase
}

func NewPaymentControllerImpl(
	presenter paymentPresenter.PaymentPresenter,
	getPaymentUseCase getpayment.GetPaymentUseCase,
	getPaymentStatusUseCase getpaymentstatus.GetPaymentStatusUseCase,
	updatePaymentUseCase updatepayment.UpdatePaymentUseCase,
	addPaymentUseCase addPayment.AddPaymentUseCase) *PaymentControllerImpl {
	return &PaymentControllerImpl{
		presenter:               presenter,
		getPaymentUseCase:       getPaymentUseCase,
		getPaymentStatusUseCase: getPaymentStatusUseCase,
		updatePaymentUseCase:    updatePaymentUseCase,
		addPaymentUseCase:       addPaymentUseCase,
	}
}

func (c *PaymentControllerImpl) CreatePayment(addPaymentRequest *dto.AddPaymentRequestDto) (string, error) {
	qrCode, err := c.addPaymentUseCase.Execute(
		commands.NewAddPaymentCommand(
			addPaymentRequest.OrderId,
			addPaymentRequest.Total,
			addPaymentRequest.Type))
	if err != nil {
		return "", err
	}

	return qrCode, nil
}

func (c *PaymentControllerImpl) GetPaymentStatusByOrderId(orderId uint) (string, error) {
	status, err := c.getPaymentStatusUseCase.Execute(commands.NewGetPaymentStatusCommand(orderId))
	if err != nil {
		return "", err
	}

	return status, nil
}

func (c *PaymentControllerImpl) GetPaymentByOrderId(orderId uint) (*dto.GetPaymentResponseDto, error) {
	payment, err := c.getPaymentUseCase.Execute(commands.NewGetPaymentCommand(orderId))
	if err != nil {
		return nil, err
	}

	return c.presenter.Present(payment), nil
}

func (c *PaymentControllerImpl) UpdatePaymentStatus(orderId uint, status string) error {
	err := c.updatePaymentUseCase.Execute(commands.NewUpdatePaymentStatusCommand(orderId, status))
	if err != nil {
		return err
	}
	return nil
}
