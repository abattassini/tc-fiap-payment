package gateways

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	paymentGateways "github.com/abattassini/tc-fiap-payment/internal/payment/gateways"
	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	"github.com/abattassini/tc-fiap-payment/pkg/rest"
)

var (
	_ paymentGateways.MercadoPagoGateway = (*MercadoPagoGatewayImpl)(nil)
)

type MercadoPagoGatewayImpl struct {
	config *MercadoPagoConfig
	client rest.HTTPClient
}

type MercadoPagoConfig struct {
	BaseURL  string
	Token    string
	Pos      string
	ClientId string
}

func (c *MercadoPagoConfig) Validate() error {
	if c.BaseURL == "" || c.Token == "" || c.Pos == "" || c.ClientId == "" {
		return fmt.Errorf("invalid MercadoPagoConfig: all fields must be set")
	}
	return nil
}

func newMercadoPagoConfig() *MercadoPagoConfig {
	return &MercadoPagoConfig{
		BaseURL:  os.Getenv("MERCADO_PAGO_BASEURL"),
		Token:    os.Getenv("MERCADO_PAGO_ACCESS_TOKEN"),
		ClientId: os.Getenv("MERCADO_PAGO_CLIENT_ID"),
		Pos:      os.Getenv("MERCADO_PAGO_POS_ID"),
	}
}

func NewMercadoPagoGatewayImpl() (*MercadoPagoGatewayImpl, error) {
	config := newMercadoPagoConfig()
	client := &http.Client{}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &MercadoPagoGatewayImpl{config: config, client: client}, nil
}

func NewMercadoPagoGatewayImplWithClient(client rest.HTTPClient) (*MercadoPagoGatewayImpl, error) {
	config := newMercadoPagoConfig()
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &MercadoPagoGatewayImpl{config: config, client: client}, nil
}

func (a *MercadoPagoGatewayImpl) buildRequest(ctx context.Context, request dto.CreateQRCodeDTO) (*http.Request, error) {
	url := fmt.Sprintf("%s/instore/orders/qr/seller/collectors/%s/pos/%s/qrs", a.config.BaseURL, a.config.ClientId, a.config.Pos)
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.config.Token)
	return req, nil
}

func (a *MercadoPagoGatewayImpl) handleResponse(resp *http.Response) (dto.QRCodeResponseDto, error) {
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return dto.QRCodeResponseDto{}, fmt.Errorf("failed to generate QR code, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var response dto.QRCodeResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.QRCodeResponseDto{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return response, nil
}

func (s *MercadoPagoGatewayImpl) GenerateQRCode(ctx context.Context, request dto.CreateQRCodeDTO) (dto.QRCodeResponseDto, error) {
	req, err := s.buildRequest(ctx, request)
	if err != nil {
		return dto.QRCodeResponseDto{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return dto.QRCodeResponseDto{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	return s.handleResponse(resp)
}
