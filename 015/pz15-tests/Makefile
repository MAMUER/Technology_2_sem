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

# Тестирование
test:
	go test ./...

test-v:
	go test -v ./...

test-internal:
	go test -v ./internal/...

test-cover:
	go test -cover ./...

test-cover-internal:
	go test -cover ./internal/...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

test-coverage-internal:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out

test-coverage-internal-html:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out

# Пакетные тесты
test-math:
	go test -v -cover ./internal/mathx/...

test-strings:
	go test -v -cover ./internal/stringsx/...

test-service:
	go test -v -cover ./internal/service/...

# Проверка покрытия для зачёта (≥70%) - Windows версия
check-coverage:
	@echo === Checking Test Coverage (min 70%) ===
	@echo Mathx package:
	@go test -cover ./internal/mathx/... | findstr "coverage:" || echo Coverage: 0.0%
	@echo Stringsx package:
	@go test -cover ./internal/stringsx/... | findstr "coverage:" || echo Coverage: 0.0%
	@echo Service package:
	@go test -cover ./internal/service/... | findstr "coverage:" || echo Coverage: 0.0%

# Альтернативная проверка покрытия (простая)
check-coverage-simple:
	@echo === Test Coverage Summary ===
	@go test -cover ./internal/mathx/...
	@go test -cover ./internal/stringsx/...
	@go test -cover ./internal/service/...

# Детальная проверка покрытия с отчётом
coverage-report:
	@echo === Detailed Coverage Report ===
	@go test -coverprofile=mathx.out ./internal/mathx/...
	@go tool cover -func=mathx.out | findstr "total:"
	@go test -coverprofile=stringsx.out ./internal/stringsx/...
	@go tool cover -func=stringsx.out | findstr "total:"
	@go test -coverprofile=service.out ./internal/service/...
	@go tool cover -func=service.out | findstr "total:"
	@del mathx.out stringsx.out service.out 2>nul

clean:
	@del bin\server.exe 2>nul
	@del coverage.out 2>nul
	@del mathx.out stringsx.out service.out 2>nul

# Бенчмарки
bench:
	go test -bench . ./internal/mathx

bench-math:
	go test -bench . ./internal/mathx -benchmem

bench-all:
	go test -bench . ./... -benchmem

bench-math-verbose:
	go test -bench . ./internal/mathx -v

help:
	@echo === Available Commands ===
	@echo BUILD:
	@echo   run           - Run the application
	@echo   build         - Build the application
	@echo   check         - Run vet and fmt
	@echo   fast          - Check, build and run
	@echo.
	@echo TESTING:
	@echo   test          - Run all tests
	@echo   test-v        - Run all tests with verbose output
	@echo   test-internal - Run internal tests with verbose output
	@echo   test-cover    - Run all tests with coverage
	@echo   test-cover-internal - Run internal tests with coverage
	@echo.
	@echo COVERAGE:
	@echo   test-coverage     - Generate coverage report (text)
	@echo   test-coverage-html - Generate coverage report (HTML)
	@echo   test-coverage-internal - Internal coverage (text)
	@echo   test-coverage-internal-html - Internal coverage (HTML)
	@echo   check-coverage    - Quick coverage check for all packages
	@echo   check-coverage-simple - Simple coverage check
	@echo   coverage-report   - Detailed coverage report per package
	@echo.
	@echo PACKAGE TESTS:
	@echo   test-math     - Run math tests with coverage
	@echo   test-strings  - Run string tests with coverage
	@echo   test-service  - Run service tests with coverage
	@echo.
	@echo UTILITY:
	@echo   clean         - Clean build artifacts and coverage files
	@echo   tree          - Show project structure
	@echo   help          - Show this help message
	@echo.
	@echo BENCHMARKS:
	@echo   bench          - Run math benchmarks
	@echo   bench-math     - Run math benchmarks with memory stats
	@echo   bench-all      - Run all benchmarks with memory stats
	@echo   bench-math-verbose - Run math benchmarks with verbose output

.PHONY: run build check tree fast \
		test test-v test-internal test-cover test-cover-internal \
		test-coverage test-coverage-html test-coverage-internal test-coverage-internal-html \
		test-math test-strings test-service \
		check-coverage check-coverage-simple coverage-report \
		clean help