# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    BINARY := bin/server.exe
    RUN_CMD := .\$(BINARY)
    TREE_CMD := tree /f
else
    DETECTED_OS := Linux
    BINARY := bin/server
    RUN_CMD := ./$(BINARY)
    TREE_CMD := tree
endif

# Основные команды
run: $(BINARY)
	$(RUN_CMD)

build:
	go build -o $(BINARY) ./cmd/server

check:	
	go vet ./...
	go fmt ./...

help:
	@echo "Test users:"
	@echo "  admin@example.com / secret123 (role: admin)"
	@echo "  user@example.com  / secret123 (role: user)"
	@echo "  user2@example.com / secret123 (role: user)"

tree:
	@if command -v tree >/dev/null 2>&1; then \
		$(TREE_CMD); \
	else \
		echo "Tree command not found. Installing..."; \
		if command -v apt-get >/dev/null 2>&1; then \
			sudo apt-get update && sudo apt-get install -y tree; \
		elif command -v yum >/dev/null 2>&1; then \
			sudo yum install -y tree; \
		elif command -v apk >/dev/null 2>&1; then \
			sudo apk add tree; \
		else \
			echo "Cannot install tree automatically. Please install manually: sudo apt-get install tree"; \
			echo "Alternative: using find command:"; \
			find . -type f -name "*.go" | head -20; \
		fi; \
	fi

fast:
	go vet ./...
	go fmt ./...
	go build -o $(BINARY) ./cmd/server
	$(RUN_CMD)

# Docker команды
docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v

docker-restart: docker-down docker-up

# Тестовые команды
test-login-admin:
	curl -X POST http://localhost:8084/api/v1/login \
	  -H "Content-Type: application/json" \
	  -d '{"email":"admin@example.com","password":"secret123"}'

test-login-user:
	curl -X POST http://localhost:8084/api/v1/login \
	  -H "Content-Type: application/json" \
	  -d '{"email":"user@example.com","password":"secret123"}'

test-health:
	curl http://localhost:8084/health

.PHONY: run build docker-up docker-down docker-logs docker-clean \
        docker-restart check test-login-admin test-login-user \
        test-health help tree fast