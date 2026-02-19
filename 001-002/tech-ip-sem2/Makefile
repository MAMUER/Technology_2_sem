ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    RUN_AUTH_CMD := cd services/auth && set AUTH_PORT=8081&& set AUTH_GRPC_PORT=50051&& go run ./cmd/auth
    RUN_TASKS_CMD := cd services/tasks && set TASKS_PORT=8082&& set AUTH_GRPC_ADDR=193.233.175.221:50051&& go run ./cmd/tasks
    TREE_CMD := tree /f
    DOCKER_COMPOSE_CMD := docker compose
else
    DETECTED_OS := Linux
    RUN_AUTH_CMD := cd services/auth && AUTH_PORT=8081 AUTH_GRPC_PORT=50051 go run ./cmd/auth
    RUN_TASKS_CMD := cd services/tasks && TASKS_PORT=8082 AUTH_GRPC_ADDR=193.233.175.221:50051 go run ./cmd/tasks
    TREE_CMD := tree
    DOCKER_COMPOSE_CMD := docker compose
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

fast-auth: check
	@echo Quick start Auth service on HTTP port $(AUTH_PORT) and gRPC port $(AUTH_GRPC_PORT)...
	$(RUN_AUTH_CMD)

fast-tasks: check
	@echo Quick start Tasks service on port $(TASKS_PORT)...
	@echo Using Auth gRPC at $(AUTH_GRPC_ADDR)
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
	@echo Auth service:  http://$(SERVER_IP):8081
	@echo Tasks service: http://$(SERVER_IP):8082
	@echo ==================================================
	@echo Use 'make docker-logs' to see logs
	@echo Use 'make docker-down' to stop all services

# Generate gRPC code
generate:
	@echo "Generating gRPC code..."
	@mkdir -p $(PROTO_OUT)
	protoc --proto_path=$(PROTO_PATH) \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "Done!"

# Install protoc plugins if needed
install-protoc:
	$(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go
	$(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc

# Utils
tree:
	@$(TREE_CMD)

help:
	@echo PORTS:
	@echo   Auth service:    $(AUTH_PORT)
	@echo   Tasks service:   $(TASKS_PORT)
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
	@echo UTILS:
	@echo   make check        - go mod tidy, vet and fmt
	@echo   make tree         - Show project structure
	@echo   make help         - This help

.PHONY: check fast-auth fast-tasks test test-docker curl-examples tree help
.PHONY: docker-build docker-up docker-down docker-restart docker-run docker-logs docker-logs-auth docker-logs-tasks docker-ps docker-clean docker-auth-shell docker-tasks-shell