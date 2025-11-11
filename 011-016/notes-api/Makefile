# Определить ОС
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    BINARY := bin/server.exe
    RUN_CMD := .\$(BINARY)
    TREE_CMD := tree /f
    KILL_CMD := taskkill /F /IM
    FIND_PROCESS := tasklist | findstr
    FIND_CMD := findstr
    SLEEP_CMD := timeout
else
    DETECTED_OS := Linux
    BINARY := bin/server
    RUN_CMD := ./$(BINARY)
    TREE_CMD := tree
    KILL_CMD := pkill -f
    FIND_PROCESS := pgrep -f
    FIND_CMD := grep
    SLEEP_CMD := sleep
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

# Быстрый запуск в фоне (кросс-платформенный)
ifeq ($(OS),Windows_NT)
fast:
	go vet ./...
	go fmt ./...
	go build -o $(BINARY) ./cmd/api
	@echo Starting server in background on port 8080...
	@start /B $(RUN_CMD)
	@echo Server started in background! Use 'make stop' to stop it.
	@echo Test endpoints:
	@echo   http://localhost:8080/debug/pprof/
	@echo   http://localhost:8080/api/v1/notes
	@echo   http://localhost:8080/work-slow
else
fast:
	go vet ./...
	go fmt ./...
	go build -o $(BINARY) ./cmd/api
	@echo "Starting server in background on port 8080..."
	@$(RUN_CMD) &
	@echo "Server PID: $$!"
	@$(SLEEP_CMD) 2
	@echo "Server started! Use 'make stop' to stop it."
	@echo "Test endpoints:"
	@echo "  http://localhost:8080/debug/pprof/"
	@echo "  http://localhost:8080/api/v1/notes"
	@echo "  http://localhost:8080/work-slow"
endif

# Остановка сервера
stop:
	@echo Stopping server...
	@$(KILL_CMD) $(BINARY) 2>nul || echo Process not found or already stopped.
	@echo Server stopped.

# Проверка статуса сервера
status:
ifeq ($(OS),Windows_NT)
	@tasklist | findstr $(BINARY) >nul && echo Server is running. || echo Server is not running.
else
	@if $(FIND_PROCESS) $(BINARY) > /dev/null; then \
		echo "Server is running"; \
	else \
		echo "Server is not running"; \
	fi
endif

# Docker команды
docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-logs-app:
	docker-compose logs app -f

docker-logs-postgres:
	docker-compose logs postgres -f

docker-status:
	docker-compose ps

docker-stats:
	docker stats

docker-clean:
	docker-compose down -v
	docker system prune -f

docker-restart:
	docker-compose restart

# PostgreSQL команды
postgres-shell:
	docker exec -it postgres-pz11 psql -U teacher_app -d notes

# Профилирование (запускать при работающем сервере)
profile-cpu:
	go tool pprof -http=:9999 http://localhost:8085/debug/pprof/profile

profile-heap:
	go tool pprof -http=:9998 http://localhost:8085/debug/pprof/heap

profile-block:
	go tool pprof -http=:9997 http://localhost:8085/debug/pprof/block

profile-mutex:
	go tool pprof -http=:9996 http://localhost:8085/debug/pprof/mutex

profile-goroutine:
	go tool pprof -http=:9995 http://localhost:8085/debug/pprof/goroutine

# Нагрузочное тестирование
load-slow:
	hey -n 100 -c 10 http://localhost:8085/work-slow

load-fast:
	hey -n 200 -c 20 http://localhost:8085/work-fast

load-notes:
	hey -n 300 -c 30 http://localhost:8085/api/v1/notes

load-pprof:
	@echo "=== Testing PPROF endpoints ==="
	hey -n 50 -c 5 http://localhost:8085/debug/pprof/
	hey -n 30 -c 3 http://localhost:8085/debug/pprof/heap
	hey -n 20 -c 2 http://localhost:8085/debug/pprof/goroutine

# Комплексное нагрузочное тестирование
load-test-all:
	@echo "=== Comprehensive Load Testing ==="
	@echo "1. Testing slow endpoint..."
	hey -n 50 -c 5 http://localhost:8085/work-slow
	@echo "2. Testing fast endpoint..."
	hey -n 100 -c 10 http://localhost:8085/work-fast
	@echo "3. Testing notes API..."
	hey -n 200 -c 20 http://localhost:8085/api/v1/notes

# Мониторинг
monitor:
	@echo "=== Real-time Monitoring ==="
	@echo "Health:"
	@curl -s http://localhost:8085/health
	@echo ""
	@echo "Docker Stats:"
	@docker stats --no-stream app-pz11 postgres-pz11

# Бенчмарки
bench:
	go test -bench=. -benchmem ./...

bench-math:
	go test -bench . ./internal/mathx -benchmem

bench-strings:
	go test -bench . ./internal/stringsx -benchmem

bench-all:
	go test -bench . ./... -benchmem

# Тестирование
test:
	go test ./...

test-v:
	go test -v ./...

test-cover:
	go test -cover ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Интеграционное тестирование
integration-test:
	go test -v ./integration/... -timeout=5m

# Утилиты
clean:
	rm -rf bin/ coverage.out

help:
	@echo "=== Available Commands (Notes API - Profiling) ==="
	@echo "BUILD:"
	@echo "  run           - Run the application"
	@echo "  build         - Build the application"
	@echo "  check         - Run vet and fmt"
	@echo "  fast          - Check, build and run in background"
	@echo "  stop          - Stop background server"
	@echo "  status        - Check server status"
	@echo ""
	@echo "DOCKER:"
	@echo "  docker-up     - Start Docker containers"
	@echo "  docker-down   - Stop Docker containers"
	@echo "  docker-logs   - Show all logs"
	@echo "  docker-logs-app - Show app logs"
	@echo "  docker-logs-postgres - Show PostgreSQL logs"
	@echo "  docker-status - Container status"
	@echo "  docker-stats  - Resource statistics"
	@echo "  docker-clean  - Clean Docker resources"
	@echo "  postgres-shell - Connect to PostgreSQL"
	@echo ""
	@echo "PROFILING:"
	@echo "  profile-cpu      - CPU profile (web UI :9999)"
	@echo "  profile-heap     - Heap profile (web UI :9998)"
	@echo "  profile-block    - Block profile (web UI :9997)"
	@echo "  profile-mutex    - Mutex profile (web UI :9996)"
	@echo "  profile-goroutine - Goroutine profile (web UI :9995)"
	@echo ""
	@echo "LOAD TESTING:"
	@echo "  load-slow      - Test slow endpoint"
	@echo "  load-fast      - Test fast endpoint"
	@echo "  load-notes     - Test notes API"
	@echo "  load-pprof     - Test pprof endpoints"
	@echo "  load-test-all  - Comprehensive load test"
	@echo "  monitor        - Real-time monitoring"
	@echo ""
	@echo "BENCHMARKS:"
	@echo "  bench          - Run all benchmarks"
	@echo "  bench-math     - Run math benchmarks"
	@echo "  bench-strings  - Run string benchmarks"
	@echo "  bench-all      - Run all benchmarks with memory"
	@echo ""
	@echo "TESTING:"
	@echo "  test           - Run all tests"
	@echo "  test-v         - Run tests with verbose output"
	@echo "  test-cover     - Run tests with coverage"
	@echo "  test-coverage  - Generate coverage report"
	@echo "  test-coverage-html - HTML coverage report"
	@echo "  integration-test - Run integration tests"
	@echo ""
	@echo "UTILITY:"
	@echo "  clean          - Clean build artifacts"
	@echo "  tree           - Show project structure"
	@echo "  help           - Show this help"

.PHONY: run build check tree fast stop status \
        docker-up docker-down docker-logs docker-logs-app docker-logs-postgres docker-status docker-stats docker-clean docker-restart postgres-shell \
        profile-cpu profile-heap profile-block profile-mutex profile-goroutine \
        load-slow load-fast load-notes load-pprof load-test-all monitor \
        bench bench-math bench-strings bench-all \
        test test-v test-cover test-coverage test-coverage-html integration-test \
        clean help