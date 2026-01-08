package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/abattassini/tc-fiap-payment/internal/payment/infrastructure/api/dto"
	"github.com/abattassini/tc-fiap-payment/pkg/rest"
)

type OrderClient interface {
	GetOrder(orderId uint) (*dto.OrderResponseDto, error)
	UpdateOrderStatus(orderId uint, status int) error
}

type OrderClientImpl struct {
	httpClient rest.HTTPClient
	baseURL    string
}

func NewOrderClient(httpClient rest.HTTPClient) *OrderClientImpl {
	baseURL := os.Getenv("ORDER_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}
	return &OrderClientImpl{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *OrderClientImpl) GetOrder(orderId uint) (*dto.OrderResponseDto, error) {
	url := fmt.Sprintf("%s/v1/order/%d", c.baseURL, orderId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get order: status %d, body: %s", resp.StatusCode, string(body))
	}

	var order dto.OrderResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (c *OrderClientImpl) UpdateOrderStatus(orderId uint, status int) error {
	url := fmt.Sprintf("%s/v1/order/%d/status", c.baseURL, orderId)

	// Create request body with status
	requestBody := map[string]uint{"status": uint(status)}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update order status: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
