package addpayment_test

import (
	"context"
	"errors"
	"testing"

	"github.com/abattassini/tc-fiap-payment/internal/payment/domain/entities"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	addpayment "github.com/abattassini/tc-fiap-payment/internal/payment/usecase/addPayment"
	"github.com/abattassini/tc-fiap-payment/internal/payment/usecase/commands"
	mockRepositories "github.com/abattassini/tc-fiap-payment/mocks/payment/domain/repositories"
	mockGateways "github.com/abattassini/tc-fiap-payment/mocks/payment/gateways"
	mockClients "github.com/abattassini/tc-fiap-payment/mocks/payment/infrastructure/clients"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AddPaymentUseCaseTestSuite struct {
	suite.Suite
	mockRepository  *mockRepositories.MockPaymentRepository
	mockGateway     *mockGateways.MockMercadoPagoGateway
	mockOrderClient *mockClients.MockOrderClient
	useCase         addpayment.AddPaymentUseCase
}

func (suite *AddPaymentUseCaseTestSuite) SetupTest() {
	suite.mockRepository = mockRepositories.NewMockPaymentRepository(suite.T())
	suite.mockGateway = mockGateways.NewMockMercadoPagoGateway(suite.T())
	suite.mockOrderClient = mockClients.NewMockOrderClient(suite.T())
	suite.useCase = addpayment.NewAddPaymentUseCaseImpl(
		suite.mockGateway,
		suite.mockOrderClient,
		suite.mockRepository,
	)
}

func TestAddPaymentUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AddPaymentUseCaseTestSuite))
}

func (suite *AddPaymentUseCaseTestSuite) Test_AddPayment_WithValidData_ShouldGenerateQRCode() {
	// GIVEN a valid payment command
	command := &commands.AddPaymentCommand{
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
	}

	savedPayment := &entities.Payment{
		ID:      1,
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
		Status:  "pending",
	}

	order := &dto.OrderResponseDto{
		ID:          1,
		TotalAmount: 100.50,
		Products: []*dto.OrderProductDto{
			{
				ProductId:   1,
				Name:        "Produto Teste",
				Description: "Descrição Teste",
				Category:    1,
				Price:       50.25,
				Quantity:    2,
			},
		},
	}

	qrCodeResponse := dto.QRCodeResponseDto{
		QRData: "00020101021243650016COM.MERCADOLIBRE",
	}

	suite.mockRepository.EXPECT().
		AddPayment(mock.MatchedBy(func(p *entities.Payment) bool {
			return p.OrderId == 1 && p.Total == 100.50 && p.Status == "pending"
		})).
		Return(savedPayment, nil).
		Once()

	suite.mockOrderClient.EXPECT().
		GetOrder(uint(1)).
		Return(order, nil).
		Once()

	suite.mockGateway.EXPECT().
		GenerateQRCode(mock.Anything, mock.MatchedBy(func(qr dto.CreateQRCodeDTO) bool {
			return qr.ExternalReference == "order-1" && qr.TotalAmount == 100.50
		})).
		Return(qrCodeResponse, nil).
		Once()

	// WHEN adding payment
	qrCode, err := suite.useCase.Execute(command)

	// THEN QR code should be generated successfully
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "00020101021243650016COM.MERCADOLIBRE", qrCode)
	suite.mockRepository.AssertExpectations(suite.T())
	suite.mockOrderClient.AssertExpectations(suite.T())
	suite.mockGateway.AssertExpectations(suite.T())
}

func (suite *AddPaymentUseCaseTestSuite) Test_AddPayment_WithRepositoryError_ShouldReturnError() {
	// GIVEN a valid payment command
	command := &commands.AddPaymentCommand{
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
	}

	expectedError := errors.New("database error")

	suite.mockRepository.EXPECT().
		AddPayment(mock.Anything).
		Return(nil, expectedError).
		Once()

	// WHEN adding payment with repository error
	qrCode, err := suite.useCase.Execute(command)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	assert.Empty(suite.T(), qrCode)
	suite.mockRepository.AssertExpectations(suite.T())
}

func (suite *AddPaymentUseCaseTestSuite) Test_AddPayment_WithOrderClientError_ShouldReturnError() {
	// GIVEN a valid payment command
	command := &commands.AddPaymentCommand{
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
	}

	savedPayment := &entities.Payment{
		ID:      1,
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
		Status:  "pending",
	}

	expectedError := errors.New("order service unavailable")

	suite.mockRepository.EXPECT().
		AddPayment(mock.Anything).
		Return(savedPayment, nil).
		Once()

	suite.mockOrderClient.EXPECT().
		GetOrder(uint(1)).
		Return(nil, expectedError).
		Once()

	// WHEN order client fails
	qrCode, err := suite.useCase.Execute(command)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	assert.Empty(suite.T(), qrCode)
	suite.mockRepository.AssertExpectations(suite.T())
	suite.mockOrderClient.AssertExpectations(suite.T())
}

func (suite *AddPaymentUseCaseTestSuite) Test_AddPayment_WithMercadoPagoError_ShouldReturnError() {
	// GIVEN a valid payment command
	command := &commands.AddPaymentCommand{
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
	}

	savedPayment := &entities.Payment{
		ID:      1,
		OrderId: 1,
		Total:   100.50,
		Type:    "QRCode",
		Status:  "pending",
	}

	order := &dto.OrderResponseDto{
		ID:          1,
		TotalAmount: 100.50,
		Products: []*dto.OrderProductDto{
			{
				ProductId:   1,
				Name:        "Produto Teste",
				Description: "Descrição Teste",
				Category:    1,
				Price:       50.25,
				Quantity:    2,
			},
		},
	}

	expectedError := errors.New("mercado pago gateway error")

	suite.mockRepository.EXPECT().
		AddPayment(mock.Anything).
		Return(savedPayment, nil).
		Once()

	suite.mockOrderClient.EXPECT().
		GetOrder(uint(1)).
		Return(order, nil).
		Once()

	suite.mockGateway.EXPECT().
		GenerateQRCode(context.Background(), mock.Anything).
		Return(dto.QRCodeResponseDto{}, expectedError).
		Once()

	// WHEN mercado pago fails
	qrCode, err := suite.useCase.Execute(command)

	// THEN error should be returned
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	assert.Empty(suite.T(), qrCode)
	suite.mockRepository.AssertExpectations(suite.T())
	suite.mockOrderClient.AssertExpectations(suite.T())
	suite.mockGateway.AssertExpectations(suite.T())
}
