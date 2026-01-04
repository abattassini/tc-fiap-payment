# Multi-stage build for smaller image size
# Stage 1: Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (often needed for Go modules)
RUN apk add --no-cache git ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first (for better layer caching)
COPY go.mod go.sum ./

# Download dependencies (cached if go.mod/go.sum haven't changed)
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/api

# Stage 2: Final minimal image
FROM scratch

# Copy ca-certificates from builder (needed for HTTPS requests)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from builder stage
COPY --from=builder /app/main /main

# Expose the application port
EXPOSE 8082

# Command to run the application
ENTRYPOINT ["/main"]
