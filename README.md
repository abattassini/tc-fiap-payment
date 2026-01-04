# TC-FIAP Payment Microservice

Payment microservice extracted from the tc-fiap-50 monolith.

## Overview

This microservice handles all payment-related operations including:
- Create payments for orders
- Get payment status
- Update payment status
- Handle payment gateway webhooks

## Environment Variables

- `DB_HOST` - Database host
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name (default: payment_db)
- `DB_SSLMODE` - SSL mode (default: disable)
- `PORT` - Application port (default: 8082)

## Running Locally

### Using Docker Compose

```bash
docker-compose up
```

The service will be available at `http://localhost:8082`

### Using Go

```bash
go mod download
go run cmd/api/main.go
```

## Development

### Run tests

```bash
make test
```

### Generate test coverage

```bash
make coverage
```

### Generate mocks

```bash
make mocks
```
