#!/bin/bash

# Load environment variables from .env file
set -a
source .env
set +a

# Start the application
cd "$(dirname "$0")"
go run cmd/api/main.go
