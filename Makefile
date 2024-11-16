build: ## Build the project
	@echo "Building the project..."
	go build -o main cmd/random-chat/main.go

start: ## Run the project
	go run cmd/random-chat/main.go

test: ## Run all tests
	go test ./...
	
install: ## Install dependencies
	go mod tidy

swag: ##  Generate swagger specification
	swag init -g cmd/random-chat/main.go

help: ## Display help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' Makefile
