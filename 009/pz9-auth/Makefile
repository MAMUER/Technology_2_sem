# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    BINARY := bin/server.exe
    RUN_CMD := .\$(BINARY)
    TREE_CMD := tree /f
    KILL_CMD := taskkill /F /IM
    FIND_PROCESS := tasklist | findstr
    FIND_PORT := netstat -an | findstr
    PSQL_CMD := psql
else
    DETECTED_OS := Linux
    BINARY := bin/server
    RUN_CMD := ./$(BINARY)
    TREE_CMD := tree
    KILL_CMD := pkill -f
    FIND_PROCESS := pgrep -f
    FIND_PORT := netstat -an | grep
    PSQL_CMD := psql
endif

# Основные команды
run: $(BINARY)
	$(RUN_CMD)

build:
	go build -o $(BINARY) ./cmd/api

check:	
	go vet ./...
	go fmt ./...

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

# Остановка туннелей
tunnel-stop:
ifeq ($(DETECTED_OS),Windows)
	$(KILL_CMD) ssh.exe 2>nul || echo No SSH tunnels found
	$(KILL_CMD) cmd.exe 2>nul || echo No CMD windows found
else
	$(KILL_CMD) ssh 2>/dev/null || echo No SSH tunnels found
endif

# Проверка подключения
check-db-root:
	$(PSQL_CMD) postgres://postgres:1@localhost:5433/todo?sslmode=disable -c "SELECT version();"

check-db:
	$(PSQL_CMD) postgres://teacher_app:secure_password_123@localhost:5433/todo?sslmode=disable -c "SELECT version();"

# Проверка подключения к контейнеру
check-db-docker:
	docker exec postgres-pz9 psql -U teacher_app -d todo -c "SELECT version();"

# Статус туннелей
tunnel-status:
ifeq ($(DETECTED_OS),Windows)
	@echo === Active SSH Processes ===
	@$(FIND_PROCESS) ssh || echo No active SSH processes found
	@echo === Listening Ports ===
	@$(FIND_PORT) 5433 || echo Port 5433 is not listening
	@$(FIND_PORT) 8083 || echo Port 8083 is not listening
else
	@echo === Active SSH Processes ===
	@ps aux | grep ssh || echo No active SSH processes found
	@echo === Listening Ports ===
	@$(FIND_PORT) 5433 || echo Port 5433 is not listening
	@$(FIND_PORT) 8083 || echo Port 8083 is not listening
endif

# Статус Docker контейнеров
docker-status:
	docker-compose ps

# Разработка
dev: check build
	$(RUN_CMD)

# Миграции (если будут нужны)
migrate: build
	$(RUN_CMD) -migrate

.PHONY: run build docker-up docker-down docker-logs docker-clean docker-restart \
        tunnel-stop check-db check-db-root check-db-docker tunnel-status \
        docker-status check test tree fast dev migrate