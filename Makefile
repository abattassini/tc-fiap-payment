.PHONY: help test test-short coverage coverage-report mocks mocks-clean mocks-regenerate build run docker-up docker-down docker-logs lint fmt vet deps deps-tidy deps-verify clean clean-all dev test-all ci

APP_NAME=tc-fiap-payment
MAIN_PATH=./cmd/api
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

ifeq ($(OS),Windows_NT)
	BINARY_EXT=.exe
	RM=cmd /C del /Q /F
	RMDIR=cmd /C rd /S /Q
	OPEN=start
else
	BINARY_EXT=
	RM=rm -f
	RMDIR=rm -rf
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		OPEN=open
	else
		OPEN=xdg-open
	endif
endif

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  help                 Show this help message"
	@echo "  test                 Run all tests"
	@echo "  test-short           Run tests without verbose output"
	@echo "  coverage             Run tests with coverage"
	@echo "  coverage-report      Generate and open coverage report"
	@echo "  mocks                Generate mocks using mockery"
	@echo "  mocks-clean          Clean generated mocks"
	@echo "  mocks-regenerate     Clean and regenerate all mocks"
	@echo "  build                Build the application"
	@echo "  run                  Run the application"
	@echo "  docker-up            Start Docker services"
	@echo "  docker-down          Stop Docker services"
	@echo "  docker-logs          Show Docker logs"
	@echo "  lint                 Run linter"
	@echo "  fmt                  Format code"
	@echo "  vet                  Run go vet"
	@echo "  deps                 Download dependencies"
	@echo "  deps-tidy            Tidy dependencies"
	@echo "  deps-verify          Verify dependencies"
	@echo "  clean                Clean build artifacts and coverage files"
	@echo "  clean-all            Clean everything including mocks"
	@echo "  dev                  Start development environment"
	@echo "  test-all             Run mocks generation, tests and coverage"
	@echo "  ci                   Run CI pipeline (tidy, mocks, test, build)"

test: ## Run all tests
	@echo "Running tests..."
	go test ./internal/payment/... -v

test-short: ## Run tests without verbose output
	@echo "Running tests..."
	go test ./internal/payment/...

coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test ./internal/payment/... -coverprofile=$(COVERAGE_FILE) -coverpkg=./internal/payment/...
	@echo "Filtering coverage (excluding mocks and generated files)..."
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	go tool cover -func=$(COVERAGE_FILE) | grep -v "mocks" | grep -v "/app/"

coverage-report: coverage ## Generate and open coverage report
	@echo "Opening coverage report..."
	$(OPEN) $(COVERAGE_HTML)

mocks: ## Generate mocks using mockery
	@echo "Generating mocks..."
	mockery

mocks-clean: ## Clean generated mocks
	@echo "Cleaning mocks..."
	$(RMDIR) mocks 2>nul || true

mocks-regenerate: mocks-clean mocks ## Clean and regenerate all mocks

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME)$(BINARY_EXT) $(MAIN_PATH)
	@echo "Build complete: bin/$(APP_NAME)$(BINARY_EXT)"

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	go run $(MAIN_PATH)/main.go

docker-up: ## Start Docker services
	@echo "Starting Docker services..."
	docker compose up -d

docker-down: ## Stop Docker services
	@echo "Stopping Docker services..."
	docker compose down

docker-logs: ## Show Docker logs
	docker compose logs -f

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

deps-tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	go mod tidy

deps-verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	go mod verify

clean: ## Clean build artifacts and coverage files
	@echo "Cleaning build artifacts..."
	$(RM) $(COVERAGE_FILE) $(COVERAGE_HTML)
	$(RMDIR) bin 2>nul || true

clean-all: clean mocks-clean ## Clean everything including mocks

dev: docker-up run ## Start development environment

test-all: mocks-regenerate test coverage ## Run mocks generation, tests and coverage

ci: deps-tidy mocks-regenerate test build ## Run CI pipeline (tidy, mocks, test, build)
