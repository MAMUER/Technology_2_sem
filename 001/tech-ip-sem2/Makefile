ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    RUN_AUTH_CMD := cd services/auth && set AUTH_PORT=8081&& go run ./cmd/auth
    RUN_TASKS_CMD := cd services/tasks && set TASKS_PORT=8082&& set AUTH_BASE_URL=http://localhost:8081&& go run ./cmd/tasks
    TREE_CMD := tree /f
    DOCKER_COMPOSE_CMD := docker compose
else
    DETECTED_OS := Linux
    RUN_AUTH_CMD := cd services/auth && AUTH_PORT=8081 go run ./cmd/auth
    RUN_TASKS_CMD := cd services/tasks && TASKS_PORT=8082 AUTH_BASE_URL=http://localhost:8081 go run ./cmd/tasks
    TREE_CMD := tree
    DOCKER_COMPOSE_CMD := docker compose
endif

AUTH_PORT = 8081
TASKS_PORT = 8082
AUTH_BASE_URL = http://localhost:$(AUTH_PORT)

# Basic commands
check:
	go mod tidy
	go vet ./...
	go fmt ./...

fast-auth: check
	@echo Quick start Auth service on port $(AUTH_PORT)...
	$(RUN_AUTH_CMD)

fast-tasks: check
	@echo Quick start Tasks service on port $(TASKS_PORT)...
	$(RUN_TASKS_CMD)

# Docker commands
docker-build:
	@echo Building Docker images...
	$(DOCKER_COMPOSE_CMD) build

docker-up:
	@echo Starting Docker containers...
	$(DOCKER_COMPOSE_CMD) up -d
	@echo Containers started. Use 'make docker-logs' to see logs

docker-down:
	@echo Stopping Docker containers...
	$(DOCKER_COMPOSE_CMD) down

docker-restart: docker-down docker-up
	@echo Docker containers restarted

docker-logs:
	$(DOCKER_COMPOSE_CMD) logs -f

docker-logs-auth:
	$(DOCKER_COMPOSE_CMD) logs -f auth

docker-logs-tasks:
	$(DOCKER_COMPOSE_CMD) logs -f tasks

docker-ps:
	$(DOCKER_COMPOSE_CMD) ps

docker-clean:
	@echo Removing all containers, networks and volumes...
	$(DOCKER_COMPOSE_CMD) down -v

docker-auth-shell:
	$(DOCKER_COMPOSE_CMD) exec auth sh

docker-tasks-shell:
	$(DOCKER_COMPOSE_CMD) exec tasks sh

# Build and run everything with Docker
docker-run: docker-build docker-up
	@echo ==================================================
	@echo All services are running!
	@echo Auth service:  http://localhost:8081
	@echo Tasks service: http://localhost:8082
	@echo ==================================================
	@echo Use 'make docker-logs' to see logs
	@echo Use 'make docker-down' to stop all services

# Test API (works with both local and Docker runs)
test:
	@echo 1. Get token:
	@curl -s -X POST http://localhost:$(AUTH_PORT)/v1/auth/login -H "Content-Type: application/json" -H "X-Request-ID: test-001" -d "{\"username\":\"student\",\"password\":\"student\"}"
	@echo.
	@echo 2. Verify token:
	@curl -s -X GET http://localhost:$(AUTH_PORT)/v1/auth/verify -H "Authorization: Bearer demo-token-for-student" -H "X-Request-ID: test-002"
	@echo.
	@echo 3. Create task:
	@curl -s -X POST http://localhost:$(TASKS_PORT)/v1/tasks -H "Content-Type: application/json" -H "Authorization: Bearer demo-token-for-student" -H "X-Request-ID: test-003" -d "{\"title\":\"Do PZ17\",\"description\":\"split services\",\"due_date\":\"2026-01-10\"}"
	@echo.
	@echo 4. List tasks:
	@curl -s -X GET http://localhost:$(TASKS_PORT)/v1/tasks -H "Authorization: Bearer demo-token-for-student" -H "X-Request-ID: test-004"

test-docker:
	@echo Testing Docker services...
	@echo 1. Check if auth service is healthy:
	@curl -s -o /dev/null -w "Auth health: %%{http_code}\n" http://localhost:8081/v1/auth/verify -H "Authorization: Bearer demo-token-for-student"
	@echo 2. Check if tasks service is healthy:
	@curl -s -o /dev/null -w "Tasks health: %%{http_code}\n" http://localhost:8082/v1/tasks -H "Authorization: Bearer demo-token-for-student"

# Examples
curl-examples:
	@echo 1. Get token:
	@echo curl -X POST http://localhost:$(AUTH_PORT)/v1/auth/login -H "Content-Type: application/json" -H "X-Request-ID: req-001" -d "{\"username\":\"student\",\"password\":\"student\"}"
	@echo.
	@echo 2. Create task:
	@echo curl -X POST http://localhost:$(TASKS_PORT)/v1/tasks -H "Content-Type: application/json" -H "Authorization: Bearer demo-token-for-student" -H "X-Request-ID: req-002" -d "{\"title\":\"Do PZ17\",\"description\":\"split services\",\"due_date\":\"2026-01-10\"}"
	@echo.
	@echo 3. List tasks:
	@echo curl -X GET http://localhost:$(TASKS_PORT)/v1/tasks -H "Authorization: Bearer demo-token-for-student" -H "X-Request-ID: req-003"

# Utils
tree:
	@$(TREE_CMD)

help:
	@echo ==================================================
	@echo AVAILABLE COMMANDS
	@echo ==================================================
	@echo PORTS:
	@echo   Auth service:    $(AUTH_PORT)
	@echo   Tasks service:   $(TASKS_PORT)
	@echo ==================================================
	@echo.
	@echo LOCAL RUN (without Docker):
	@echo   make fast-auth    - Auth with code check
	@echo   make fast-tasks   - Tasks with code check
	@echo.
	@echo DOCKER COMMANDS:
	@echo   make docker-build     - Build Docker images
	@echo   make docker-up        - Start containers in background
	@echo   make docker-down      - Stop containers
	@echo   make docker-restart   - Restart containers
	@echo   make docker-run       - Build and start everything
	@echo   make docker-logs      - Show all logs
	@echo   make docker-logs-auth - Show auth service logs
	@echo   make docker-logs-tasks - Show tasks service logs
	@echo   make docker-ps        - Show container status
	@echo   make docker-clean     - Remove all containers and volumes
	@echo   make docker-auth-shell - Open shell in auth container
	@echo   make docker-tasks-shell - Open shell in tasks container
	@echo.
	@echo TESTING:
	@echo   make test         - Test API (local)
	@echo   make test-docker  - Test API (docker)
	@echo   make curl-examples - Show curl examples
	@echo.
	@echo UTILS:
	@echo   make check        - go mod tidy, vet and fmt
	@echo   make tree         - Show project structure
	@echo   make help         - This help
	@echo ==================================================

.PHONY: check fast-auth fast-tasks test test-docker curl-examples tree help
.PHONY: docker-build docker-up docker-down docker-restart docker-run docker-logs docker-logs-auth docker-logs-tasks docker-ps docker-clean docker-auth-shell docker-tasks-shell