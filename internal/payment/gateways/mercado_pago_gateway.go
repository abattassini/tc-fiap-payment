package gateways

import (
	"context"

	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
)

type MercadoPagoGateway interface {
	GenerateQRCode(ctx context.Context, request dto.CreateQRCodeDTO) (dto.QRCodeResponseDto, error)
}
