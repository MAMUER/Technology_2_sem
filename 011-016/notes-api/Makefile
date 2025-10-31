# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    BINARY := bin/server.exe
    RUN_CMD := .\$(BINARY)
    TREE_CMD := tree /f
    PSQL_CMD := psql
    FIND_CMD := findstr
else
    DETECTED_OS := Linux
    BINARY := bin/server
    RUN_CMD := ./$(BINARY)
    TREE_CMD := tree
    PSQL_CMD := psql
    FIND_CMD := grep
endif

# Основные команды
run: $(BINARY)
	$(RUN_CMD)

build:
	go build -o $(BINARY) ./cmd/api

check:	
	go vet ./...
	go fmt ./...

swagger:
	swag init -g cmd/api/main.go -o docs

tree:
	$(TREE_CMD)

fast:
	go vet ./...
	go fmt ./...
	go build -o $(BINARY) ./cmd/api
	$(RUN_CMD)

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

# Нагрузочное тестирование
load-test: load-test-pagination load-test-batch load-test-single load-test-create

load-test-pagination:
	@echo "=== Testing Pagination (1000 requests, 50 concurrent) ==="
	hey -n 1000 -c 50 "http://193.233.175.221:8085/api/v1/notes?limit=20"

load-test-batch:
	@echo "=== Testing Batch Requests (500 requests, 20 concurrent) ==="
	hey -n 500 -c 20 "http://193.233.175.221:8085/api/v1/notes?ids=1,2,3,4,5"

load-test-single:
	@echo "=== Testing Single Requests (2000 requests, 100 concurrent) ==="
	hey -n 2000 -c 100 "http://193.233.175.221:8085/api/v1/notes/1"

load-test-create:
	@echo "=== Testing Create Operations (300 requests, 10 concurrent) ==="
	hey -n 300 -c 10 -m POST -d "{\"title\":\"Test\",\"content\":\"Content\"}" "http://193.233.175.221:8085/api/v1/notes"

# Проверка подключения
check-db:
	$(PSQL_CMD) postgres://teacher_app:secure_password_123@193.233.175.221:5434/notes?sslmode=disable -c "SELECT version();"

check-db-docker:
	docker exec postgres-pz11 psql -U teacher_app -d notes -c "SELECT version();"

# Тестовые команды
test-create:
	curl -X POST http://193.233.175.221:8085/api/v1/notes \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer test-token" \
	  -d '{"title":"Тестовая заметка","content":"Содержание"}'

test-list:
	curl -X GET "http://193.233.175.221:8085/api/v1/notes?page=1&limit=10" \
	  -H "Authorization: Bearer test-token" \
	  -H "Content-Type: application/json"

test-get:
	curl -X GET http://193.233.175.221:8085/api/v1/notes/1 \
	  -H "Authorization: Bearer test-token" \
	  -H "Content-Type: application/json"

test-update:
	curl -X PATCH http://193.233.175.221:8085/api/v1/notes/1 \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer test-token" \
	  -d '{"title":"Обновленный заголовок"}'

test-delete:
	curl -X DELETE http://193.233.175.221:8085/api/v1/notes/1 \
	  -H "Authorization: Bearer test-token" \
	  -H "Content-Type: application/json"

test-swagger:
	curl http://193.233.175.221:8085/docs/doc.json

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

# Профили (запускать при работающем сервере)
profile-cpu:
	go tool pprof -http=:9999 http://193.233.175.221:8080/debug/pprof/profile

profile-block:
	go tool pprof -http=:9998 http://193.233.175.221:8080/debug/pprof/block

profile-mutex:
	go tool pprof -http=:9997 http://193.233.175.221:8080/debug/pprof/mutex

# Нагрузка (нужен hey: go install github.com/rakyll/hey@latest)
load:
	hey -n 100 -c 10 http://193.233.175.221:8080/work-slow

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