# Product Inventory Management System Makefile

# Variables
BIDAPI_DIR = bidapi
BIDRPC_DIR = bidrpc
PROTO_DIR = bidrpc/bidrpcproto
GO_MODULE = bidfood

# Default target
.PHONY: help
help: ## Show this help message
	@echo "Product Inventory Management System"
	@echo "==================================="
	@echo "Available commands:"
	@echo ""
	@echo "  api                run bidapi service"
	@echo "  rpc                run bidrpc service"
	@echo "  proto              generate protobuf file"
	@echo "  test               run all unit tests"
	@echo "  test-api           run api tests"
	@echo "  test-rpc           run grpcurl tests"
	@echo "  deps               install dependencies"
	@echo ""


.PHONY: api
api:
	cd $(BIDAPI_DIR) && go run ./...

.PHONY: rpc
rpc:
	cd $(BIDRPC_DIR) && go run ./...

.PHONY: proto
proto: ## Generate protobuf files
	@echo "Generating protobuf files..."
	protoc --proto_path=. \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/product.proto
	@echo "Protobuf files generated successfully!"

.PHONY: test
test:
	@echo "Running unit tests..."
	go test ./...


.PHONY: test-api
test-api:
	@echo "Running api tests..."
	./test_api.sh


.PHONY: test-rpc
test-rpc:
	@echo "Running rpc tests..."
	./test_rpc.sh

.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpcurl/cmd/grpcurl@latest
	@echo "Dependencies installed successfully!"
