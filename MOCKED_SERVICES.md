# External Services Integration Status

## Summary
This document tracks the integration status of external service dependencies.

## ✅ Integrated Services (Real API Calls)

### 1. Order Service (tc-fiap-order)
**Status**: ✅ **INTEGRATED** - Real service calls are active

Both Order Service endpoints have been integrated:

- **GetOrder()** - `internal/payment/usecase/addPayment/add_payment_use_case_impl.go`
  - **Endpoint**: `GET /v1/order/{orderId}`
  - **Purpose**: Fetches order details to create payment QR code
  - **Status**: ✅ Active

- **UpdateOrderStatus()** - `internal/payment/usecase/handleWebhook/handle_webhook_use_case_impl.go`
  - **Endpoint**: `PUT /v1/order/{orderId}/status`
  - **Request Body**: `{"status": 2}`
  - **Purpose**: Updates order status to "Preparing" when payment is approved
  - **Status**: ✅ Active

**Configuration**:
```bash
ORDER_SERVICE_URL=http://localhost:8081  # Default value
```

---

## ⚠️ Partially Mocked Services

### 2. MercadoPago Gateway
**Status**: ⚠️ **PARTIALLY MOCKED** - Uses mock mode with test credentials

- **File**: `internal/payment/infrastructure/gateways/mercado_pago_gateway_impl.go`
- **Behavior**: Automatically switches between mock/real based on credentials
  - **Mock mode**: When using `test_token_12345` or `test_client_id`
  - **Real mode**: When using valid MercadoPago credentials
- **Current**: Using TEST credentials (real API calls)

**To Use Real API**:
Replace test credentials in `.env` with production credentials.


