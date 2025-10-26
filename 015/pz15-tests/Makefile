# Основные команды
run:
	go run ./cmd/api

build:
	go build -o bin/server.exe ./cmd/api

check:	
	go vet ./...
	go fmt ./...

tree:
	tree /f

fast:
	go vet ./...
	go fmt ./...
	go build -o bin/server.exe ./cmd/api
	.\bin\server.exe

test:
	go test ./internal/... -v

test-math:
	go test ./internal/mathx/... -v

test-strings:
	go test ./internal/stringsx/... -v

test-service:
	go test ./internal/service/... -v

test-coverage:
	go test ./internal/... -cover

clean:
	rm -f bin/server.exe

help:
	@echo "Available commands:"
	@echo "  run           - Run the application"
	@echo "  build         - Build the application"
	@echo "  check         - Run vet and fmt"
	@echo "  test          - Run all tests"
	@echo "  test-math     - Run math tests"
	@echo "  test-strings  - Run string tests"
	@echo "  test-service  - Run service tests"
	@echo "  clean         - Clean build artifacts"

.PHONY: run build check tree fast test test-math test-strings test-service test-coverage clean help