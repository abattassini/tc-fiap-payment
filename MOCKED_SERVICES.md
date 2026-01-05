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

### 2. MercadoPago Gateway
- **Status**: ✅ **WORKING** - Using real TEST credentials
- **File**: `internal/payment/infrastructure/gateways/mercado_pago_gateway_impl.go`
- **Note**: Automatically switches between mock/real based on credentials
  - Mock mode: When using `test_token_12345` or `test_client_id`
  - Real mode: When using valid TEST credentials (currently active)

**Search for**: `// TODO: Replace mock` to find all mock locations

