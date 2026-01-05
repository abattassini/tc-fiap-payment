# Mocked External Services

## Summary
External services are mocked to enable local testing without dependencies.

## Mocked Locations

### 1. Order Service (2 places)
- **File**: `internal/payment/usecase/addPayment/add_payment_use_case_impl.go:48`
  - **Mock**: Returns fake order data
  - **TODO**: Replace with `u.orderClient.GetOrder()`

- **File**: `internal/payment/usecase/handleWebhook/handle_webhook_use_case_impl.go:45`
  - **Mock**: Skips order status update
  - **TODO**: Replace with `u.orderClient.UpdateOrderStatus()`

### 2. MercadoPago Gateway (1 place)
- **File**: `internal/payment/infrastructure/gateways/mercado_pago_gateway_impl.go:75`
  - **Mock**: Returns fake QR code
  - **TODO**: Uncomment real API call

**Search for**: `// TODO: Replace mock` to find all mock locations
