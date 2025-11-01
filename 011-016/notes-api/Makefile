# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    BINARY := bin/server.exe
    RUN_CMD := .\$(BINARY)
    TREE_CMD := tree /f
    KILL_CMD := taskkill /F /IM
    FIND_PROCESS := tasklist | findstr
else
    DETECTED_OS := Linux
    BINARY := bin/server
    RUN_CMD := ./$(BINARY)
    TREE_CMD := tree
    KILL_CMD := pkill -f
    FIND_PROCESS := pgrep -f
endif

# Основные команды
run:
	$(RUN_CMD)

build:
	go build -o $(BINARY) ./cmd/api

check:	
	go vet ./...
	go fmt ./...

tree:
	@$(TREE_CMD)

# Быстрый запуск в фоне
fast:
	go vet ./...
	go fmt ./...
	go build -o $(BINARY) ./cmd/api
	@echo "Starting server in background on port 8080..."
	@$(RUN_CMD) &
	@echo "Server PID: $$!"
	@sleep 2
	@echo "Server started! Use 'make stop' to stop it."
	@echo "Test endpoints:"
	@echo "  http://localhost:8080/debug/pprof/"
	@echo "  http://localhost:8080/api/v1/notes"
	@echo "  http://localhost:8080/work-slow"

# Остановка сервера
stop:
	@echo "Stopping server..."
	@$(KILL_CMD) $(BINARY) 2>/dev/null || true
	@echo "Server stopped."

# Проверка статуса сервера
status:
	@if $(FIND_PROCESS) $(BINARY) > /dev/null; then \
		echo "Server is running"; \
	else \
		echo "Server is not running"; \
	fi

# Запуск через Docker (рекомендуется для тестирования)
docker-run:
	docker-compose up -d
	@sleep 3
	@echo "Server running on http://localhost:8085"
	@echo "Test PPROF: http://localhost:8085/debug/pprof/"

# Остальные команды остаются без изменений...
swagger:
	swag init -g cmd/api/main.go -o docs

# Docker команды
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v

docker-restart: docker-down docker-up

# Профили (запускать при работающем сервере)
profile-cpu:
	go tool pprof -http=:9999 http://193.233.175.221:8085/debug/pprof/profile

profile-block:
	go tool pprof -http=:9998 http://193.233.175.221:8085/debug/pprof/block

profile-mutex:
	go tool pprof -http=:9997 http://193.233.175.221:8085/debug/pprof/mutex

# Нагрузка (нужен hey: go install github.com/rakyll/hey@latest)
load:
	hey -n 100 -c 10 http://193.233.175.221:8085/work-slow

# Нагрузочное тестирование PPROF эндпоинтов
load-pprof:
	@echo "=== Testing PPROF endpoints ==="
	hey -n 50 -c 5 http://193.233.175.221:8085/debug/pprof/
	hey -n 30 -c 3 http://193.233.175.221:8085/debug/pprof/heap
	hey -n 20 -c 2 http://193.233.175.221:8085/debug/pprof/goroutine

# Тестирование производительности
load-performance:
	@echo "=== Testing Performance Endpoints ==="
	hey -n 100 -c 10 http://193.233.175.221:8085/work-slow
	hey -n 200 -c 20 http://193.233.175.221:8085/work-fast
	hey -n 50 -c 5 http://193.233.175.221:8085/block-demo

# Нагрузочное тестирование
load-test: load-test-pagination load-test-batch load-test-single load-test-create

load-test-pagination:
	@echo "=== Testing Pagination (1000 requests, 50 concurrent) ==="
	hey -n 1000 -c 50 "http://localhost:8085/api/v1/notes?limit=20"

load-test-batch:
	@echo "=== Testing Batch Requests (500 requests, 20 concurrent) ==="
	hey -n 500 -c 20 "http://localhost:8085/api/v1/notes?ids=1,2,3,4,5"

load-test-single:
	@echo "=== Testing Single Requests (2000 requests, 100 concurrent) ==="
	hey -n 2000 -c 100 "http://localhost:8085/api/v1/notes/1"

load-test-create:
	@echo "=== Testing Create Operations (300 requests, 10 concurrent) ==="
	hey -n 300 -c 10 -m POST -d "{\"title\":\"Test\",\"content\":\"Content\"}" "http://localhost:8085/api/v1/notes"

# Бенчмарки
bench:
	go test -bench=. -benchmem ./...

bench-mthx:
	go test -bench . ./internal/mathx

bench-math:
	go test -bench . ./internal/mathx -benchmem

bench-all:
	go test -bench . ./... -benchmem

bench-math-verbose:
	go test -bench . ./internal/mathx -v

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

# Проверка покрытия для зачёта (≥70%)
check-coverage:
	@echo "=== Checking Test Coverage (min 70%) ==="
	@echo "Mathx package:"
	@go test -cover ./internal/mathx/... | $(FIND_CMD) "coverage:" || echo "Coverage: 0.0%"
	@echo "Stringsx package:"
	@go test -cover ./internal/stringsx/... | $(FIND_CMD) "coverage:" || echo "Coverage: 0.0%"
	@echo "Service package:"
	@go test -cover ./internal/service/... | $(FIND_CMD) "coverage:" || echo "Coverage: 0.0%"

# Альтернативная проверка покрытия (простая)
check-coverage-simple:
	@echo "=== Test Coverage Summary ==="
	@go test -cover ./internal/mathx/...
	@go test -cover ./internal/stringsx/...
	@go test -cover ./internal/service/...

# Детальная проверка покрытия с отчётом
coverage-report:
	@echo "=== Detailed Coverage Report ==="
	@go test -coverprofile=mathx.out ./internal/mathx/...
	@go tool cover -func=mathx.out | $(FIND_CMD) "total:"
	@go test -coverprofile=stringsx.out ./internal/stringsx/...
	@go tool cover -func=stringsx.out | $(FIND_CMD) "total:"
	@go test -coverprofile=service.out ./internal/service/...
	@go tool cover -func=service.out | $(FIND_CMD) "total:"
	@rm -f mathx.out stringsx.out service.out 2>/dev/null || true

clean:
	@rm -f $(BINARY) 2>/dev/null || true
	@rm -f coverage.out 2>/dev/null || true
	@rm -f mathx.out stringsx.out service.out 2>/dev/null || true

# Интеграционное тестирование
integration-up:
	docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 5

integration-down:
	docker-compose down -v

# Интеграционное тестирование с Testcontainers
integration-test:
	@echo "Running integration tests with Testcontainers..."
	go test -v ./integration/... -timeout=5m

integration-test-cover:
	@echo "Running integration tests with coverage..."
	go test -v ./integration/... -cover -coverprofile=integration.coverage.out -timeout=5m
	go tool cover -html=integration.coverage.out -o integration.coverage.html

# Все тесты
test-all: test integration-test

# Очистка
clean-integration:
	rm -f integration.coverage.out integration.coverage.html

help:
	@echo "=== Available Commands ==="
	@echo "BUILD:"
	@echo "  run           - Run the application"
	@echo "  build         - Build the application"
	@echo "  check         - Run vet and fmt"
	@echo "  fast          - Check, build and run"
	@echo ""
	@echo "TESTING:"
	@echo "  test          - Run all tests"
	@echo "  test-v        - Run all tests with verbose output"
	@echo "  test-internal - Run internal tests with verbose output"
	@echo "  test-cover    - Run all tests with coverage"
	@echo "  test-cover-internal - Run internal tests with coverage"
	@echo ""
	@echo "COVERAGE:"
	@echo "  test-coverage     - Generate coverage report (text)"
	@echo "  test-coverage-html - Generate coverage report (HTML)"
	@echo "  test-coverage-internal - Internal coverage (text)"
	@echo "  test-coverage-internal-html - Internal coverage (HTML)"
	@echo "  check-coverage    - Quick coverage check for all packages"
	@echo "  check-coverage-simple - Simple coverage check"
	@echo "  coverage-report   - Detailed coverage report per package"
	@echo ""
	@echo "PACKAGE TESTS:"
	@echo "  test-math     - Run math tests with coverage"
	@echo "  test-strings  - Run string tests with coverage"
	@echo "  test-service  - Run service tests with coverage"
	@echo ""
	@echo "UTILITY:"
	@echo "  clean         - Clean build artifacts and coverage files"
	@echo "  tree          - Show project structure"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "BENCHMARKS:"
	@echo "  bench          - Run math benchmarks"
	@echo "  bench-math     - Run math benchmarks with memory stats"
	@echo "  bench-all      - Run all benchmarks with memory stats"
	@echo "  bench-math-verbose - Run math benchmarks with verbose output"

.PHONY: run build check swagger tree fast \
        docker-up docker-down docker-logs docker-clean docker-restart \
        load-test load-test-pagination load-test-batch load-test-single load-test-create \
        check-db check-db-docker test-create test-list test-get test-update test-delete test-swagger \
		test test-v test-internal test-cover test-cover-internal \
		test-coverage test-coverage-html test-coverage-internal test-coverage-internal-html \
		test-math test-strings test-service \
		check-coverage check-coverage-simple coverage-report \
		clean help \
		integration-test integration-test-cover test-all clean-integration