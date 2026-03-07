ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    RUN_AUTH_CMD := cd services/auth && set AUTH_PORT=8081&& set AUTH_GRPC_PORT=50051&& go run ./cmd/auth
    RUN_TASKS_CMD := cd services/tasks && set TASKS_PORT=8082&& set AUTH_GRPC_ADDR=193.233.175.221:50051&& go run ./cmd/tasks
    TREE_CMD := tree /f
    DOCKER_COMPOSE_CMD := docker compose
    PROTO_CMD := protoc --proto_path=proto --go_out=proto/gen/go --go_opt=paths=source_relative --go-grpc_out=proto/gen/go --go-grpc_opt=paths=source_relative proto/auth.proto
else
    DETECTED_OS := Linux
    RUN_AUTH_CMD := cd services/auth && AUTH_PORT=8081 AUTH_GRPC_PORT=50051 go run ./cmd/auth
    RUN_TASKS_CMD := cd services/tasks && TASKS_PORT=8082 AUTH_GRPC_ADDR=193.233.175.221:50051 go run ./cmd/tasks
    TREE_CMD := tree
    DOCKER_COMPOSE_CMD := docker compose
    PROTO_CMD := protoc --proto_path=$(PROTO_PATH) --go_out=$(PROTO_OUT) --go_opt=paths=source_relative --go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative $(PROTO_FILES)
endif

AUTH_PORT = 8081
TASKS_PORT = 8082
AUTH_BASE_URL = http://193.233.175.221:$(AUTH_PORT)
SERVER_IP = 193.233.175.221

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

# gRPC
PROTOC=protoc
PROTO_PATH=proto
PROTO_OUT=proto/gen/go
PROTO_FILES=$(wildcard $(PROTO_PATH)/*.proto)

AUTH_GRPC_PORT = 50051
AUTH_GRPC_ADDR = localhost:$(AUTH_GRPC_PORT)

# Basic commands
check:
	go mod tidy
	go vet ./...
	go fmt ./...
	go test ./...
up:
	@echo "=== ЗАПУСК ВСЕГО ПРОЕКТА (HTTP + gRPC + Prometheus + Grafana + HTTPS + PostgreSQL) ==="
	$(MAKE) gen-cert
	$(MAKE) docker-up
	@echo ""
	@echo "=== ПРОЕКТ ЗАПУЩЕН ==="
	@echo "HTTP endpoints:"
	@echo "  Auth HTTP:  http://localhost:8081"
	@echo "  Tasks HTTP: http://localhost:8082"
	@echo "HTTPS endpoints:"
	@echo "  Auth HTTPS:  https://localhost:8443/v1/auth"
	@echo "  Tasks HTTPS: https://localhost:8443/v1/tasks"
	@echo "Monitoring:"
	@echo "  Prometheus: http://localhost:9090"
	@echo "  Grafana:    http://localhost:3000"
	@echo "Database:"
	@echo "  PostgreSQL: localhost:5432"

down:
	@echo "=== ОСТАНОВКА ВСЕГО ПРОЕКТА ==="
	$(DOCKER_COMPOSE_CMD) down
	cd deploy/tls && $(DOCKER_COMPOSE_CMD) down

restart:
	$(MAKE) down
	$(MAKE) up

# Docker commands
DOCKER_COMPOSE_CMD := docker-compose

docker-build:
	@echo Building Docker images...
	$(DOCKER_COMPOSE_CMD) build
	cd deploy/tls && $(DOCKER_COMPOSE_CMD) build

docker-up:
	@echo Starting all containers...
	$(DOCKER_COMPOSE_CMD) up -d
	# Убираем запуск отдельного стека из deploy/tls
	# cd deploy/tls && $(DOCKER_COMPOSE_CMD) up -d
	@echo Waiting for services to be ready...
	@sleep 10
	$(MAKE) docker-ps

docker-down:
	@echo Stopping all containers...
	$(DOCKER_COMPOSE_CMD) down
	cd deploy/tls && $(DOCKER_COMPOSE_CMD) down

docker-restart: docker-down docker-up
	@echo Docker containers restarted

docker-logs:
	$(DOCKER_COMPOSE_CMD) logs -f

docker-logs-all:
	$(DOCKER_COMPOSE_CMD) logs -f &
	cd deploy/tls && $(DOCKER_COMPOSE_CMD) logs -f

docker-ps:
	@echo "=== MAIN STACK ==="
	$(DOCKER_COMPOSE_CMD) ps
	@echo ""
	@echo "=== TLS STACK ==="
	cd deploy/tls && $(DOCKER_COMPOSE_CMD) ps

docker-clean:
	@echo Removing all containers and volumes...
	$(DOCKER_COMPOSE_CMD) down -v
	cd deploy/tls && $(DOCKER_COMPOSE_CMD) down -v

docker-reset:
	@echo "========== ПОЛНЫЙ СБРОС DOCKER =========="
	@echo "1. Проверка свободного места на диске..."
	@df -h .
	@echo ""
	@echo "2. Останавливаем и удаляем все контейнеры, сети и тома..."
	$(DOCKER_COMPOSE_CMD) down -v
	@echo ""
	@echo "3. Агрессивная очистка Docker (удаляем всё неиспользуемое)..."
	@docker system prune -a -f --volumes
	@docker container prune -f
	@docker image prune -a -f
	@docker volume prune -f
	@docker network prune -f
	@echo ""
	@echo "4. Проверка свободного места после очистки..."
	@df -h .
	@echo ""
	@echo "5. Собираем образы..."
	$(DOCKER_COMPOSE_CMD) build
	@echo ""
	@echo "6. Запускаем новые контейнеры..."
	$(DOCKER_COMPOSE_CMD) up -d
	@echo ""
	@echo "7. Статус контейнеров:"
	$(DOCKER_COMPOSE_CMD) ps
	@echo "=========================================="
	@echo "Полный сброс завершен! Используйте 'make docker-logs' для просмотра логов"

docker-stop-auth:
	@echo Stopping auth service container...
	docker stop pz20-auth
	@echo Auth service stopped

docker-start-auth:
	@echo Starting auth service container...
	docker start pz20-auth
	@echo Auth service started

docker-restart-auth: docker-stop-auth docker-start-auth
	@echo Auth service restarted

docker-logs-auth:
	$(DOCKER_COMPOSE_CMD) logs -f auth

docker-logs-tasks:
	$(DOCKER_COMPOSE_CMD) logs -f tasks

docker-auth-shell:
	$(DOCKER_COMPOSE_CMD) exec auth sh

docker-tasks-shell:
	$(DOCKER_COMPOSE_CMD) exec tasks sh

docker-logs-clean:
	@echo "Cleaning auth logs..."
	@sudo truncate -s 0 $$(docker inspect --format='{{.LogPath}}' pz20-auth 2>/dev/null) 2>/dev/null || echo "Auth container not running"
	@echo "Cleaning tasks logs..."
	@sudo truncate -s 0 $$(docker inspect --format='{{.LogPath}}' pz20-tasks 2>/dev/null) 2>/dev/null || echo "Tasks container not running"
	@echo "Logs cleared!"

docker-logs-color:
	$(DOCKER_COMPOSE_CMD) logs -f --tail=100 | grep --color=always -E "error|warn|info|debug|$$"

docker-health:
	@echo "Auth service health:"
	@curl -s -o /dev/null -w "%{http_code}" http://localhost:8081/v1/auth/verify -H "Authorization: Bearer demo-token-for-student" || echo "Failed"
	@echo ""
	@echo "Tasks service health:"
	@curl -s -o /dev/null -w "%{http_code}" http://localhost:8082/v1/tasks -H "Authorization: Bearer demo-token-for-student" || echo "Failed"
	@echo ""

# Generate gRPC
generate:
	@echo "Generating gRPC code..."
ifeq ($(OS),Windows_NT)
	@if not exist proto\gen\go\auth mkdir proto\gen\go\auth
	@C:\protoc\bin\protoc --proto_path=proto \
		--go_out=. --go_opt=module=tech-ip-sem2 \
		--go-grpc_out=. --go-grpc_opt=module=tech-ip-sem2 \
		proto/auth.proto
else
	@mkdir -p proto/gen/go/auth
	@protoc --proto_path=proto \
		--go_out=. --go_opt=module=tech-ip-sem2 \
		--go-grpc_out=. --go-grpc_opt=module=tech-ip-sem2 \
		proto/auth.proto
endif
	@echo "Done!"

# Monitoring commands
monitor-up:
	@echo Starting monitoring stack...
	docker-compose up -d prometheus grafana

monitor-down:
	@echo Stopping monitoring stack...
	docker-compose stop prometheus grafana

monitor-logs:
	docker-compose logs -f prometheus grafana

monitor-ps:
	docker-compose ps prometheus grafana

monitor-restart:
	docker-compose restart prometheus grafana

# Generate test load
test-load:
	@echo "Generating test load..."
	@for i in {1..50}; do \
		curl -s -X GET http://localhost:8082/v1/tasks \
			-H "Authorization: Bearer demo-token-for-student" \
			-H "X-Request-ID: load-test-$$i" > /dev/null & \
	done
	@echo "50 successful requests sent"
	@sleep 2
	@for i in {1..20}; do \
		curl -s -X GET http://localhost:8082/v1/tasks \
			-H "Authorization: Bearer wrong-token" \
			-H "X-Request-ID: load-test-error-$$i" > /dev/null & \
	done
	@echo "20 error requests sent"

test-load-create:
	@echo "Creating test tasks..."
	@for i in {1..10}; do \
		curl -s -X POST http://localhost:8082/v1/tasks \
			-H "Content-Type: application/json" \
			-H "Authorization: Bearer demo-token-for-student" \
			-H "X-Request-ID: create-$$i" \
			-d "{\"title\":\"Test Task $$i\",\"description\":\"Test description\",\"due_date\":\"2026-02-20\"}" > /dev/null & \
	done
	@echo "10 tasks created"

# Generate self-signed certificates for HTTPS
gen-cert:
	@echo "Generating self-signed certificates..."
	@mkdir -p deploy/tls/certs
ifeq ($(OS),Windows_NT)
	@echo "========================================================"
	@echo "Windows detected - please run this command on server:"
	@echo "========================================================"
	@echo "ssh root@193.233.175.221"
	@echo "cd /opt/techip"
	@echo "mkdir -p certs"
	@echo "openssl req -x509 -newkey rsa:2048 -nodes -keyout certs/key.pem -out certs/cert.pem -days 365 -subj /CN=localhost"
	@echo "docker-compose restart nginx"
	@echo "========================================================"
else
	@openssl req -x509 -newkey rsa:2048 -nodes \
		-keyout deploy/tls/certs/key.pem \
		-out deploy/tls/certs/cert.pem \
		-days 365 \
		-subj "/CN=localhost"
	@chmod 644 deploy/tls/certs/key.pem deploy/tls/certs/cert.pem
	@echo "Certificates generated in deploy/tls/certs/"
endif

p22-up:
	@echo "Starting HTTPS + PostgreSQL stack..."
	cd deploy/tls && docker-compose up -d  

p22-down:
	@echo "Stopping HTTPS + PostgreSQL stack..."
	cd deploy/tls && docker-compose down

p22-logs:
	cd deploy/tls && docker-compose logs -f

p22-ps:
	cd deploy/tls && docker-compose ps

p22-restart:
	cd deploy/tls && docker-compose restart

p22-clean:
	cd deploy/tls && docker-compose down -v

# Practice 6 - CSRF and XSS Protection
p23-up:
	@echo "Starting CSRF/XSS protection demo..."
	$(MAKE) gen-cert
	$(MAKE) docker-up

p23-show-cookies:
	@echo "Current cookies:"
	@cat cookies.txt 2>/dev/null || echo "No cookies file found"

# Redis commands
redis-up:
	@echo "Starting Redis cluster..."
	cd deploy/redis && docker-compose up -d

redis-down:
	@echo "Stopping Redis cluster..."
	cd deploy/redis && docker-compose down

redis-logs:
	cd deploy/redis && docker-compose logs -f

redis-ps:
	cd deploy/redis && docker-compose ps

redis-cli:
	docker exec -it pz20-redis redis-cli

redis-flush:
	docker exec pz20-redis redis-cli FLUSHALL
	@echo "Redis cache flushed"

redis-info:
	docker exec pz20-redis redis-cli INFO

# Practice 9 - Redis Cache
p25-up:
	@echo "Starting Redis cache demo..."
	$(MAKE) gen-cert
	$(MAKE) docker-up

p25-test-hit:
	@echo "Testing cache hit (first request - should be MISS)..."
	curl -s -X GET http://localhost:8082/v1/tasks \
		-H "Authorization: Bearer demo-token-for-student" \
		-H "X-Request-ID: cache-test-1" | jq .
	@echo ""
	@echo "Second request - should be HIT..."
	curl -s -X GET http://localhost:8082/v1/tasks \
		-H "Authorization: Bearer demo-token-for-student" \
		-H "X-Request-ID: cache-test-2" | jq .

p25-test-degrade:
	@echo "Stopping Redis to test degradation..."
	docker stop pz20-redis
	@echo "Making request without Redis (should work)..."
	curl -s -X GET http://localhost:8082/v1/tasks \
		-H "Authorization: Bearer demo-token-for-student" \
		-H "X-Request-ID: degrade-test" | jq .
	@echo ""
	@echo "Restarting Redis..."
	docker start pz20-redis

p25-clear-cache:
	@echo "Clearing Redis cache..."
	docker exec pz20-redis redis-cli FLUSHALL
	@echo "Cache cleared"

# Practice 10 - Load Balancing
LB_DIR = deploy/lb

lb-up:
	@echo "Starting load balancing demo with 3 replicas..."
	cd $(LB_DIR) && docker-compose up -d --build
	@echo ""
	@echo "=== LOAD BALANCER READY ==="
	@echo "Load Balancer: http://localhost:8080"
	@echo "Replicas: tasks_1, tasks_2, tasks_3"
	@echo ""
	@echo "Test commands:"
	@echo "  make lb-test-roundrobin  - Test round-robin distribution"
	@echo "  make lb-test-health      - Check health of all replicas"
	@echo "  make lb-test-failover    - Test failover (stop one replica)"

lb-down:
	@echo "Stopping load balancing demo..."
	cd $(LB_DIR) && docker-compose down

lb-logs:
	cd $(LB_DIR) && docker-compose logs -f

lb-ps:
	cd $(LB_DIR) && docker-compose ps

lb-restart:
	cd $(LB_DIR) && docker-compose restart

lb-clean:
	cd $(LB_DIR) && docker-compose down -v

# Test round-robin distribution (10 requests)
lb-test-roundrobin:
	@echo "Testing round-robin distribution (10 requests)..."
	@for i in {1..10}; do \
		echo "Request $$i: "; \
		curl -s -I http://localhost:8080/v1/tasks \
			-H "Authorization: Bearer demo-token-for-student" \
			-H "X-Request-ID: lb-test-$$i" | grep -i "X-Instance-ID"; \
	done

# Test with JSON response
lb-test-roundrobin-json:
	@echo "Testing round-robin distribution with JSON output..."
	@for i in {1..10}; do \
		INSTANCE=$$(curl -s -D - http://localhost:8080/v1/tasks \
			-H "Authorization: Bearer demo-token-for-student" \
			-H "X-Request-ID: lb-test-$$i" | grep -i "X-Instance-ID" | tr -d '\r'); \
		echo "Request $$i: $$INSTANCE"; \
	done

# Check health of all replicas
lb-test-health:
	@echo "Checking health of all replicas..."
	@echo "--- tasks_1 ---"
	@docker exec pz26-tasks-1 curl -s http://localhost:8082/health | jq .
	@echo ""
	@echo "--- tasks_2 ---"
	@docker exec pz26-tasks-2 curl -s http://localhost:8082/health | jq .
	@echo ""
	@echo "--- tasks_3 ---"
	@docker exec pz26-tasks-3 curl -s http://localhost:8082/health | jq .

# Test failover (stop one replica)
lb-test-failover:
	@echo "=== TESTING FAILOVER ==="
	@echo "Initial state - all 3 replicas running:"
	@make lb-test-roundrobin | head -3
	@echo ""
	@echo "Stopping tasks_2..."
	@docker stop pz26-tasks-2
	@sleep 2
	@echo ""
	@echo "After stopping tasks_2 - requests should go only to tasks_1 and tasks_3:"
	@make lb-test-roundrobin | head -5
	@echo ""
	@echo "Restarting tasks_2..."
	@docker start pz26-tasks-2
	@sleep 2
	@echo ""
	@echo "After restart - all 3 replicas back:"
	@make lb-test-roundrobin | head -3

# Test with different HTTP methods
lb-test-methods:
	@echo "Testing different HTTP methods through load balancer..."
	@echo ""
	@echo "GET /health:"
	curl -s http://localhost:8080/health | jq .
	@echo ""
	@echo "GET /v1/tasks (list):"
	curl -s -X GET http://localhost:8080/v1/tasks \
		-H "Authorization: Bearer demo-token-for-student" \
		-H "X-Request-ID: lb-method-1" | jq .
	@echo ""
	@echo "POST /v1/tasks (create):"
	curl -s -X POST http://localhost:8080/v1/tasks \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer demo-token-for-student" \
		-H "X-Request-ID: lb-method-2" \
		-d '{"title":"LB Test","description":"Through load balancer","due_date":"2026-03-15"}' | jq .

# Test search through load balancer
lb-test-search:
	@echo "Testing search through load balancer..."
	curl -s "http://localhost:8080/v1/tasks/search?q=Test" \
		-H "Authorization: Bearer demo-token-for-student" \
		-H "X-Request-ID: lb-search" | jq .

# Test metrics through load balancer
lb-test-metrics:
	@echo "Getting metrics through load balancer..."
	curl -s http://localhost:8080/metrics | grep -E "http_requests_total|http_request_duration" | head -10

# Clean up and restart everything
lb-reset: lb-clean lb-up
	@echo "Load balancer demo reset complete"

# Utils
tree:
	@$(TREE_CMD)

help:
	@echo === MAIN COMMANDS ===
	@echo   make up          
	@echo   make down        
	@echo   make restart     
	@echo   make test        
	@echo 
	@echo PORTS:
	@echo   Auth service HTTP: $(AUTH_PORT)
	@echo   Auth service gRPC: $(AUTH_GRPC_PORT)
	@echo   Tasks service:     $(TASKS_PORT)
	@echo.
	@echo LOCAL RUN (without Docker):
	@echo   make generate      - Generate gRPC code from proto
	@echo   make fast-auth     - Auth with code check
	@echo   make fast-tasks    - Tasks with code check
	@echo.
	@echo DOCKER COMMANDS:
	@echo   make docker-build     - Build Docker images
	@echo   make docker-up        - Start containers in background
	@echo   make docker-down      - Stop containers
	@echo   make docker-restart   - Restart containers
	@echo   make docker-reset     - Full reset: down, clean, rebuild, up
	@echo   make docker-reset-fast - Quick reset without cache cleanup
	@echo   make docker-logs      - Show all logs
	@echo   make docker-logs-auth - Show auth service logs
	@echo   make docker-logs-tasks - Show tasks service logs
	@echo   make docker-logs-color - Show logs with color highlighting
	@echo   make docker-ps        - Show container status
	@echo   make docker-clean     - Remove all containers and volumes
	@echo   make docker-auth-shell - Open shell in auth container
	@echo   make docker-tasks-shell - Open shell in tasks container
	@echo   make docker-stop-auth   - Stop only auth service
	@echo   make docker-start-auth  - Start only auth service
	@echo   make docker-restart-auth - Restart only auth service
	@echo   make docker-logs-clean  - Clear container logs
	@echo   make docker-health      - Check service health
	@echo.
	@echo UTILS:
	@echo   make check        - go mod tidy, vet and fmt
	@echo   make tree         - Show project structure
	@echo   make help         - This help

.PHONY: check fast-auth fast-tasks test test-docker curl-examples tree help generate gen-cert
.PHONY: docker-build docker-up docker-down docker-restart docker-reset docker-reset-fast
.PHONY: docker-logs docker-logs-auth docker-logs-tasks docker-logs-color docker-ps docker-clean
.PHONY: docker-auth-shell docker-tasks-shell docker-stop-auth docker-start-auth docker-restart-auth
.PHONY: docker-logs-clean docker-health